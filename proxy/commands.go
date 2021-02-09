package proxy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/ldelossa/vim-grpc.vim/channel"
	pb "github.com/ldelossa/vim-grpc.vim/proto/commands"
)

// CommandRecord is a record structure for book-keeping
// registered extenion commands.
type CommandRecord struct {
	registration *pb.CommandRegistration
	stream       pb.Commands_RegisterCommandServer
}

// CommandsService handles extension command book-keeping (registration, listing, deleting)
// and monitors the channel's mailboxes for incoming Command rpcs.
type CommandsService struct {
	*Proxy
	pb.UnimplementedCommandsServer
	sync.Mutex
	cmds map[string]CommandRecord
}

func NewCommandsService(ctx context.Context, proxy *Proxy) *CommandsService {
	cs := &CommandsService{
		Proxy: proxy,
		cmds:  map[string]CommandRecord{},
	}
	for i := channel.CMDBoxNumOffset; i < channel.RPCBoxNumOffset; i++ {
		go cs.monitor(ctx, i)
	}
	return cs
}

// monitor watches the provided mailbox number for incoming Command rpcs.
//
// when monitor encounters a Command rpc it will forward this event to the
// extension which registered it.
func (c *CommandsService) monitor(ctx context.Context, boxNumber uint32) {
	var cmdEvent pb.CommandIssued
	var ch channel.Channel
	for ctx.Err() == nil {
		runtime.Gosched()

		ch = c.Channel()
		if !ch.ChannelOpen() {
			continue
		}

		d := channel.Delivery{Channel: ch, BoxNum: boxNumber}
		env, err := d.Wait(ctx)
		if err != nil {
			log.Printf("CommandsService: received error waiting on mailbox %v: %v", boxNumber, err)
			continue
		}

		err = jsonpb.Unmarshal(bytes.NewReader(env.Body), &cmdEvent)
		if err != nil {
			log.Printf("CommandsService: received error serializing json to CommandIssued event %v: %v", boxNumber, err)
			continue
		}

		err = c.doCommand(&cmdEvent)
		if err != nil {
			log.Printf("CommandService: failed to do command: %v", err)
			continue
		}
	}
	log.Printf("CommandsService: monitor ctx canceled: %v", ctx.Err())
	return
}

func (c *CommandsService) doCommand(event *pb.CommandIssued) error {
	c.Lock()
	defer c.Unlock()
	cmd, ok := c.cmds[event.Command]
	if !ok {
		return fmt.Errorf("command " + event.Command + "does not exist")
	}
	if err := cmd.stream.Send(&pb.CommandEvent{Event: &pb.CommandEvent_Issued{Issued: event}}); err != nil {
		return err
	}
	return nil
}

// RegisterCommand will attempt to register the provided command with Vim.
// On success the Server side stream will be held open and the client will receive
// receipt of a command invocation via calling Recv on its side of the stream.
//
// If the channel to Vim disconnects the stream will be closed.
// If the client disconnects the stream will be closed.
// In both occurrences the registered command is removed from Vim.
func (c *CommandsService) RegisterCommand(req *pb.RegisterCommandRequest, stream pb.Commands_RegisterCommandServer) error {
	const (
		RPC = "RegisterCommand"
	)

	ch := c.Channel()
	if !ch.ChannelOpen() {
		return channel.ErrChanClosed
	}

	m := jsonpb.Marshaler{
		EmitDefaults: false,
	}

	var b bytes.Buffer
	err := m.Marshal(&b, req)
	if err != nil {
		return err
	}

	e := channel.Envelope{
		RPC:  RPC,
		Body: b.Bytes(),
	}

	e, err = ch.Send(stream.Context(), &e).Wait(stream.Context())
	if err != nil {
		return err
	}

	var cmdReg pb.CommandRegistration
	err = jsonpb.Unmarshal(bytes.NewReader(e.Body), &cmdReg)
	if err != nil {
		return fmt.Errorf("CommandsService: failed decoding expected CommandRegistration event: %v", err)
	}

	if cmdReg.Registered != true {
		return fmt.Errorf("registration failed: %v", cmdReg.Reason)
	}

	c.Lock()
	c.cmds[req.Command] = CommandRecord{
		registration: &cmdReg,
		stream:       stream,
	}
	c.Unlock()

	stream.Send(&pb.CommandEvent{
		Event: &pb.CommandEvent_Registration{Registration: &cmdReg},
	})

	<-stream.Context().Done()
	return stream.Context().Err()
}
