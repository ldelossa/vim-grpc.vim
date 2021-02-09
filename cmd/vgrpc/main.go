// EnvService
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	cmds "github.com/ldelossa/vim-grpc.vim/proto/commands"
	env "github.com/ldelossa/vim-grpc.vim/proto/env"
	"github.com/ldelossa/vim-grpc.vim/proxy"
	"google.golang.org/grpc"
)

const (
	GRPCListenAddr = "localhost:8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// create and start the proxy.
	// creates the tcp socket vim will
	// connect to.
	p := proxy.NewProxy(ctx)
	log.Printf("starting proxy on localhost:%v", proxy.DefaultPort)
	go func() {
		err := p.Listen(ctx)
		if err != nil {
			log.Printf("failed to start proxy: %v", err)
			cancel()
		}
	}()

	// create and start the gRPC server.
	// registers the gRPC server and Proxy service
	// gRPC clients will connect to.
	lis, err := net.Listen("tcp", GRPCListenAddr)
	if err != nil {
		log.Fatalf("failed to create gRPC listener: %v", err)
	}

	grpcServer := grpc.NewServer()

	env.RegisterEnvServer(grpcServer, p)
	cmds.RegisterCommandsServer(grpcServer, p)

	log.Printf("starting grpc server on %v", GRPCListenAddr)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("error starting grpc server: %v", err)
			cancel()
		}
	}()

	// block main thread on sigint or ctx cancelation.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	select {
	case <-sig:
		log.Printf("received sigint. gracefully shutting down")
	case <-ctx.Done():
		// error logged already by proxy or grpc go routine
		// if we got here.
	}
	cancel()
	grpcServer.GracefulStop()
}
