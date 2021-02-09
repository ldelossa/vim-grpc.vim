package channel

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
)

// Delivery provides a lock-free wait mechanism
// for clients issueing a channel.Send
//
// Delivery will check its Channel's mailbox periodically
// until an envelope is available.
type Delivery struct {
	Channel Channel
	BoxNum  uint32
	Err     error
}

// Wait will block until a Vim response is delivered or the ctx is canceled.
// Any errors from the call to Channel.Send() will be returned by Wait's
// err value.
func (d *Delivery) Wait(ctx context.Context) (env Envelope, err error) {
	if d.Err != nil {
		return Envelope{}, d.Err
	}
	for {
		if *d.Channel.State == Closed {
			return Envelope{}, fmt.Errorf("channel closed")
		}
		if ctx.Err() != nil {
			return Envelope{}, ctx.Err()
		}
		e := (*Envelope)(*d.Channel.mailbox[int(d.BoxNum)])
		if e != nil && e.In {
			env = *e
			atomic.SwapPointer(d.Channel.mailbox[int(env.Mailbox)], nil)
			return
		}
		runtime.Gosched()
	}
}
