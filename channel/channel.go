package channel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	Closed int32 = iota
	Open
)

var ErrChanClosed = errors.New("channel closed")

const (
	RPCBoxNumOffset uint32 = 4
	CMDBoxNumOffset uint32 = 0
)

// Channel represents a Vim channel in JSON mode.
//
// A Channel maintains a mailbox where incoming RPC messages
// are placed, in the form of Envelope data structures.
//
// Each Delivery polls its assigned mailbox until an Envelope is present.
//
// Mailbox numbers 0-3 are reserved for broadcasting registered commands.
//
// In the occurence of an underlying TCP error the Channel delivers
// a sentinel Envelope indicating error.
//
// A Channel is safe to pass by value.
type Channel struct {
	// following pointers guaranteed non nil if constructor
	// is used.
	State *int32 // atomically updated
	conn  *net.TCPConn
	*json.Encoder
	*json.Decoder
	mailbox []*unsafe.Pointer
}

// Close should be called on TCP terminating errors.
// The TCP conn will be closed and an synethic error
// RPC will be provided to all clients.
func (c *Channel) Close() {
	e := &Envelope{
		Err: fmt.Errorf("channel closed during rpc"),
	}
	// one caller will win this race.
	ok := atomic.CompareAndSwapInt32(c.State, 1, 0)
	if ok {
		c.conn.Close()
		for i, p := range c.mailbox {
			if (*Envelope)(*p) != nil && !(*Envelope)(*p).In {
				atomic.SwapPointer(c.mailbox[i], unsafe.Pointer(e))
			}
		}
	}
}

// Ping sends synethic rpcs as a
// heartbeat with Vim.
//
// Ping will close the connetion and return
// if an underlying tcp error is encountered.
// In this case Ping unblocks.
//
// Ping will also unblock if the channel enters
// a closed state or the provided ctx is canceled.
func (c Channel) Ping(ctx context.Context) {
	e := &Envelope{
		RPC: "Ping",
	}
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		if *c.State == Closed {
			log.Printf("channel closed during ping")
			return
		}
		select {
		case <-ctx.Done():
			log.Printf("channel: ctx cancled, closing channel: %v", ctx.Err())
			return
		case <-t.C:
			tctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
			env, err := c.Send(tctx, e).Wait(tctx)
			if err != nil {
				log.Printf("channel: err while sending ping, closing channel: %v", err)
				c.Close()
				cancel()
				return
			}
			if env.RPC != "Pong" {
				log.Printf("channel: ping did not receive pong: %v", env.RPC)
				c.Close()
				cancel()
				return
			}
			log.Printf("channel: received pong!")
			cancel()
		}
	}
}

func NewChannel(conn *net.TCPConn) Channel {
	open := int32(1)
	c := Channel{
		Encoder: json.NewEncoder(conn),
		Decoder: json.NewDecoder(conn),
		conn:    conn,
		mailbox: make([]*unsafe.Pointer, 1024),
		State:   &open,
	}
	// initialize unsafes
	for i := 0; i < 1024; i++ {
		p := unsafe.Pointer(nil)
		c.mailbox[i] = &p
	}
	return c
}

// Send delivers an Envelope to Vim.
//
// The Envelope's BoxNum field should be empty, and is ignored
// if not.
//
// The Envelope's In field MUST be false.
//
// Any error encountered on a Send is deferred
// until the Wait() call on the provided Deliverer.
func (c Channel) Send(ctx context.Context, e *Envelope) *Delivery {
	if *c.State == Closed {
		return &Delivery{Err: fmt.Errorf("channel closed")}
	}
	// try to compare and swap a mailbox number
	// until success or ctx cancel
	boxNum := RPCBoxNumOffset
	for ; boxNum < uint32(len(c.mailbox)); boxNum++ {
		if ctx.Err() != nil {
			return &Delivery{
				Err: ctx.Err(),
			}
		}
		e.Mailbox = boxNum
		if ok := atomic.CompareAndSwapPointer(c.mailbox[boxNum], nil, unsafe.Pointer(e)); ok {
			break
		}
	}

	vim, err := e.ToVim()
	if err != nil {
		return &Delivery{
			Err: fmt.Errorf("failed to encode to vim type: %v", err),
		}
	}

	err = c.Encode(vim)
	if err != nil {
		log.Printf("channel: error sending, closing channel: %v", err)
		c.Close()
		return &Delivery{Err: err}
	}

	return &Delivery{
		Channel: c,
		BoxNum:  boxNum,
	}
}

// Recv reads off the json.Decoder
// until the ctx is canceled or the underlying
// tcp.conn fails.
//
// When a json message is received the payload
// will be placed in the channel's mailbox.
func (c Channel) Recv(ctx context.Context) {
	for {
		if *c.State == Closed {
			log.Printf("channel: channel closed during recv")
			return
		}
		if ctx.Err() != nil {
			log.Printf("channel rcv ctx canceled: %v", ctx.Err())
			return
		}
		vw, e := VimWrap{}, &Envelope{}
		err := c.Decode(&vw)
		if err != nil {
			log.Printf("channel: error receiving, closing channel: %v", err)
			c.Close()
			break
		}
		if err = e.FromVim(vw); err != nil {
			log.Printf("channel: could not create Envelope from VimWrap: %v", err)
		}
		e.In = true
		atomic.SwapPointer(c.mailbox[int(e.Mailbox)], unsafe.Pointer(e))
	}
}

// ChannelOpen reports whether the channel is opened or not.
func (c Channel) ChannelOpen() bool {
	switch {
	case *(c.State) == Open:
		return true
	case *(c.State) == Closed:
		return false
	}
	panic("unreachable")
}
