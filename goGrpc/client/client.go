package main

import (
	"bufio"
	"context"
	"fmt"
	"grpcChatServer/chatserver/chatserver"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	//Create a variable and then use the same variable in every gorutine?

	fmt.Println("Enter Server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connecting : " + serverID)

	//Connect client to server
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server :: %v", err)
	}

	//Create a stream with ChatService
	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	//Establish communication with gRPC Server
	ch := clienthandle{stream: stream}
	// Create channel to send Lamport timestamp between different gorutines.
	bl := make(chan bool)
	ch.clientConfig()
	go ch.sendMessage(bl)
	go ch.receiveMessage()

	<-bl

	defer conn.Close() //closes the connection once everything has returned
}

// Maybe save Lamport timestamp inside clienthandle since it is already used by all gorutines?
type clienthandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clienthandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")
	//User joined Messsage
	clientMessageBox := &chatserver.FromClient{
		Name: ch.clientName,
		Body: "joined Chitty-Chat at Lamport time L",
	}
	err = ch.stream.Send(clientMessageBox)

}

// aright what is this stuff
func (ch *clienthandle) sendMessage(bl chan bool) {
	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")
		if clientMessage == "/leave" {
			//User Left Message
			clientMessageBox := &chatserver.FromClient{
				Name: ch.clientName,
				Body: "left Chitty-Chat at Lamport time L",
			}
			err = ch.stream.Send(clientMessageBox)
			bl <- false
		} else if len(clientMessage) > 128 {
			fmt.Printf("Message is above the allowed 128 characters\n")
		} else {

			clientMessageBox := &chatserver.FromClient{
				Name: ch.clientName,
				Body: clientMessage,
			}

			err = ch.stream.Send(clientMessageBox)

			if err != nil {
				log.Printf("Error while sending message to server :: %v", err)
			}
		}
	}

}

func (ch *clienthandle) receiveMessage() {
	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		//print message to console
		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)

	}
}
