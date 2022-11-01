package main

import (
	"bufio"
	"flag"
	"grpcChatServer/chatserver"
	"log"
	"os"
	"strconv"
	"sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var name string
var lamport int32 = 1
var client chatserver.ServicesClient
var context_ context.Context

func Join() {
	stream, _ := client.Join(context.Background(), &chatserver.JoinMessage{User: name})

	for {
		response, err := stream.Recv()

		if err != nil {
			break
		}

		atomic.StoreInt32(&lamport, MaxInt(lamport, response.Lamport)+1)

		if response.User == "" {
			log.Default().Printf("(%s) >> %s", strconv.Itoa(int(lamport)), response.Content)
			continue
		}

		log.Default().Printf("(%s, %s) >> %s", strconv.Itoa(int(lamport)), response.User, response.Content)
	}
}

func MaxInt(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func Publish(message string) {
	if len(message) > 128 {
		log.Print("Message is above the allowed 128 characters\n")
		return
	}

	if len(message) == 0 {
		log.Print("Message is empty")
		return
	}

	atomic.AddInt32(&lamport, 1)
	_, err := client.Publish(context_, &chatserver.Message{User: name, Content: message, Lamport: lamport})

	if err != nil {
		log.Fatalf("Error while sending message to server: %s", err)
	}
}

func main() {
	nameFlag := flag.String("name", "", "")
	flag.Parse()
	name = *nameFlag

	// Connect to server
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to server! %s", err)
		return
	}

	defer conn.Close()

	client = chatserver.NewServicesClient(conn)
	context_ = context.Background()

	go Join()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go Publish(scanner.Text())
	}
}
