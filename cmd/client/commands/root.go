package commands

import (
	"context"
	"fmt"
	"os"

	pb "github.com/ldelossa/vim-grpc.vim/proto/commands"
	"google.golang.org/grpc"
)

const (
	help = `
The 'commands' sub-command is used to register extension commands with Vim.
register - registers a command with Vim
list     - list all registered commands 
exec     - execute a command

`
)

func Root(ctx context.Context, conn *grpc.ClientConn) error {
	if len(os.Args) < 3 {
		fmt.Print(help)
		return fmt.Errorf("error: needs subcommand")
	}

	client := pb.NewCommandsClient(conn)

	sub := os.Args[2]
	switch sub {
	case "register":
		return register(ctx, client)
	default:
		return fmt.Errorf("error: unknown subcommand: %v", sub)
	}
}
