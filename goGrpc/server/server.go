package main

import (
	"grpcChatServer/chatserver/chatserver"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" //default port set to 5000
	}

	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Print("Listening at: " + Port)

	//Server instance
	grpcserver := grpc.NewServer()

	//Chatservice
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpcserver, &cs)

	//Fail to start grpc server
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}

}
