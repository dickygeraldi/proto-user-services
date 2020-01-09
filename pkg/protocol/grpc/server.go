package grcp

import (
	"context"
	"log"
	"os"
	"os/signal"

	"fmt"
	"net"
	v1 "protoUserService/pkg/api/v1"

	"google.golang.org/grpc"
)

// Run server gRPC sevice to user service
func RunServer(ctx context.Context, v1Api v1.UserServicesServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}

	// Register the service
	server := grpc.NewServer()
	v1.RegisterUserServicesServer(server, v1Api)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Shutting down gRPC server.....")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("Start gRPC server")
	return server.Serve(listen)
}
