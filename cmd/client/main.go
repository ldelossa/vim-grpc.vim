package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ldelossa/vim-grpc.vim/cmd/client/commands"
	"google.golang.org/grpc"
)

const (
	DefaultGRPCServerAddr = "localhost:8080"
	help                  = `This CLI is used to test vim-grpc functionality and is split into a series of subcommands.

The following subcommands are available:

commands - this command is used to register extension commands with vim-grpc and logs a message when the command issued at Vim.
`
)

func main() {
	conn, err := grpc.Dial(DefaultGRPCServerAddr,
		grpc.WithTimeout(5*time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure())
	if err != nil {
		fmt.Printf("error: gRPC server dial failed: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Print(help)
		fmt.Print("need at least one root command.\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "commands":
		err := commands.Root(context.TODO(), conn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
