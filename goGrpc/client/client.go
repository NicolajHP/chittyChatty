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
var client chatserver.ServicesClient
var ctx context.Context
var lamport int32 = 1

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
		log.Print("The message must not be above 128 characters!")
		return
	}

	if len(message) == 0 {
		log.Print("The message cannot be empty!")
		return
	}

	atomic.AddInt32(&lamport, 1)
	_, err := client.Publish(ctx, &chatserver.Message{User: name, Content: message, Lamport: lamport})

	if err != nil {
		log.Fatalf("Could not send the message.. Error: %s", err)
	}
}

func main() {
	// Handle flags
	nameFlag := flag.String("name", "", "")

	flag.Parse()
	name = *nameFlag

	// Handle connection
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect! %s", err)
		return
	}

	defer conn.Close()

	client = chatserver.NewServicesClient(conn)
	ctx = context.Background()

	go Join()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go Publish(scanner.Text())
	}
}
