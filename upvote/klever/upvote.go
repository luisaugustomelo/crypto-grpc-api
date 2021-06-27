package klever

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) UpVote(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %v", message.Body)

	return &Message{Body: "Hello From the Server!"}, nil
}
