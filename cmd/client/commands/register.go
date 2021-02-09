package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/ldelossa/vim-grpc.vim/proto/commands"
)

var registerFS = flag.NewFlagSet("commands register", flag.ExitOnError)

var registerFlags = struct {
	extension *string
	command   *string
	title     *string
}{
	extension: registerFS.String("ext", "", "name of the extension registering this command (required)"),
	command:   registerFS.String("cmd", "", "name of the command being registered (required)"),
	title:     registerFS.String("title", "", "title of the command in Vim (required)"),
}

func register(ctx context.Context, client pb.CommandsClient) error {
	registerFS.Usage = func() {
		fmt.Print(`Usage of commands register:
  -cmd string
        name of the command being registered (required)
  -ext string
        name of the extension registering this command (required)
  -title string
        title of the command in Vim (required)

On successful registration issuing the command at Vim will log a message here for testing.
`)
	}
	registerFS.Parse(os.Args[3:])

	if *registerFlags.extension == "" {
		return fmt.Errorf("'ext' argument required")
	}
	if *registerFlags.command == "" {
		return fmt.Errorf("'cmd' argument required")
	}
	if *registerFlags.title == "" {
		return fmt.Errorf("'title' argument required")
	}

	req := &pb.RegisterCommandRequest{
		Extension: *registerFlags.extension,
		Command:   *registerFlags.command,
		Title:     *registerFlags.title,
	}
	stream, err := client.RegisterCommand(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register command: %v", err)
	}
	defer stream.CloseSend()

	event, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("error on first recv: %v", err)
	}

	reg := event.GetRegistration()
	if reg == nil {
		return fmt.Errorf("first message was not a registration message")
	}

	if !reg.Registered {
		return fmt.Errorf("registration failed: %v", reg.Reason)
	}

	for {
		event, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("error on rcv: %v", err)
		}
		issued := event.GetIssued()
		if issued == nil {
			return fmt.Errorf("received unhandled message, closing command channel.")
		}
		log.Printf("command triggered: %+v", issued)
	}
}
