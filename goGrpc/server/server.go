package main

import (
	"grpcChatServer/chatserver"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	log.Print("Loading...")

	listener, err := net.Listen("tcp", "localhost:5000")

	if err != nil {
		log.Fatalf("Could not listen @ %s", err)
		return
	}

	log.Print("setting up server...")

	s := chatserver.Server{}

	grpcServer := grpc.NewServer()

	chatserver.RegisterServicesServer(grpcServer, &s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("Failed to start gRPC Server :: %v", err)
	}
}
