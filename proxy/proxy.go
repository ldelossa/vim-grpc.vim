package proxy

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/ldelossa/vim-grpc.vim/channel"
)

const (
	Network     = "tcp4"
	DefaultPort = 7999
)

// Proxy implements our gRPC client <-> Vim
// multiplexing proxy.
//
// The Proxy embeds individual gRPC services exposing a namespaced
// RPC API to gRPC clients.
//
// A Proxy must always be constructed by its NewProxy constructor
// to properly initialize a disconnected channel.
type Proxy struct {
	// embedded services promote gRPC service implementation methods
	// making Proxy usable in all gRPC registration
	// methods.
	*EnvironmentService
	*CommandsService
	sync.RWMutex
	channel channel.Channel
}

func NewProxy(ctx context.Context) *Proxy {
	p := &Proxy{}
	p.channel = channel.Channel{State: new(int32)}
	// register services.
	p.EnvironmentService = NewEnvService(ctx, p)
	p.CommandsService = NewCommandsService(ctx, p)
	return p
}

// Listen will create a TCP socket for
// Vim to connect to.
//
// Once a connection is made a channel.Channel
// will be created from the net.TCPConn
//
// Any gRPC requests made while a Channel is not
// initialized will error.
//
// Proxy will block on a single Vim channel
// and will not concurrently handle multiple.
func (p *Proxy) Listen(ctx context.Context) error {
	listener, err := net.ListenTCP(Network, &net.TCPAddr{
		Port: DefaultPort,
	})
	if err != nil {
		return err
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("err tcp during connection: %v", err)
			continue
		}
		log.Printf("proxy: received new connect")

		p.Lock()
		p.channel = channel.NewChannel(conn)
		p.Unlock()

		log.Printf("proxy: channel connected")
		// kick off recv side
		go p.channel.Recv(ctx)
		// blocks until ctx is canceled or an underlying
		// tcp error is detected.
		p.channel.Ping(ctx)
		log.Printf("proxy: channel disconnected")
	}
}

// Channel returns the Vim channel associated with the
// Proxy.
//
// The Channel is not guaranteed to be open.
func (p *Proxy) Channel() channel.Channel {
	p.Lock()
	ch := p.channel
	p.Unlock()
	return ch
}
