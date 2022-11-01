package chatserver

import (
	"log"
	"strconv"

	"golang.org/x/net/context"

	"sync/atomic"
)

type Server struct {
	UnimplementedServicesServer
}

var clients []Services_JoinServer = make([]Services_JoinServer, 0)
var lamport int32 = 1

func (s *Server) Broadcast(ctx context.Context, message *Message) (*Empty, error) {
	atomic.AddInt32(&lamport, 1)
	atomic.StoreInt32(&message.Lamport, lamport)

	for _, client := range clients {
		client.Send(message)
	}

	if message.User != "" {
		log.Printf("(%s, %s) >> %s", strconv.Itoa(int(message.Lamport)), message.User, message.Content)
	} else {
		log.Printf("(%s) >> %s", strconv.Itoa(int(message.Lamport)), message.Content)
	}
	return &Empty{}, nil
}

func (s *Server) Join(message *JoinMessage, stream Services_JoinServer) error {
	clients = append(clients, stream)

	msg := Message{
		User:    "",
		Content: "Participant " + message.User + " joined Chitty-Chat at Lamport time " + strconv.Itoa(int(lamport)),
		Lamport: lamport,
	}

	s.Broadcast(context.TODO(), &msg)

	for {
		select {
		case <-stream.Context().Done():
			msg := Message{
				User:    "",
				Content: "Participant " + message.User + " left Chitty-Chat at Lamport time " + strconv.Itoa(int(lamport)),
				Lamport: lamport,
			}
			for i, element := range clients {
				if element == stream {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			s.Broadcast(context.TODO(), &msg)
			return nil
		}
	}

}

func (s *Server) Publish(ctx context.Context, message *Message) (*Empty, error) {
	atomic.StoreInt32(&lamport, MaxInt(lamport, message.Lamport)+1)
	return s.Broadcast(ctx, message)
}

func MaxInt(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
