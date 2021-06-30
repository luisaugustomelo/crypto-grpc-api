package system

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

func (s *Server) CreateCrypto(ctx context.Context, message *CreateCryptocurrencyRequest) (*CreateCryptocurrencyResponse, error) {
	crypto := Cryptocurrency{Id: "1", Name: "2", Initials: "3", Description: "4", Upvote: 5, Downvote: 6}

	return &CreateCryptocurrencyResponse{Crypto: &crypto}, nil
}
