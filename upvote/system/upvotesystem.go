package system

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) PingPong(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %v", message.Body)

	return &Message{Body: "Hello From the Server!"}, nil
}

func (s *Server) CreateCryptocurrency(ctx context.Context, message *CreateCryptocurrencyRequest) (*CreateCryptocurrencyResponse, error) {
	crypto := Cryptocurrency{Id: "1", Name: "2", Initials: "3", Description: "4", Upvote: 5, Downvote: 6}

	return &CreateCryptocurrencyResponse{Crypto: &crypto}, nil
}

func (s *Server) UpdateCryptocurrency(ctx context.Context, message *UpdateCryptocurrencyRequest) (*UpdateCryptocurrencyResponse, error) {
	return &UpdateCryptocurrencyResponse{}, nil
}

func (s *Server) DeleteCryptocurrency(ctx context.Context, message *DeleteCryptocurrencyRequest) (*DeleteCryptocurrencyResponse, error) {
	return &DeleteCryptocurrencyResponse{}, nil
}

func (s *Server) ReadCryptocurrencyById(ctx context.Context, message *ReadCryptocurrencyRequest) (*ReadCryptocurrencyResponse, error) {
	return &ReadCryptocurrencyResponse{}, nil
}

func (s *Server) ListAllCriptocurrencies(ctx context.Context, message *ListAllCryptocurrenciesRequest) (*ListAllCryptocurrenciesResponse, error) {
	return &ListAllCryptocurrenciesResponse{}, nil
}

func (s *Server) UpVoteCriptocurrency(ctx context.Context, message *UpVoteCryptocurrencyRequest) (*UpVoteCryptocurrencyResponse, error) {
	return &UpVoteCryptocurrencyResponse{}, nil
}

func (s *Server) DownVoteCriptocurrency(ctx context.Context, message *DownVoteCryptocurrencyRequest) (*DownVoteCryptocurrencyResponse, error) {
	return &DownVoteCryptocurrencyResponse{}, nil
}

func (s *Server) GetSumVotes(ctx context.Context, message *GetSumVotesRequest) (*GetSumVotesResponse, error) {
	return &GetSumVotesResponse{}, nil
}

func (s *Server) GetSumVotesByStreamRequest(ctx context.Context, message *GetSumVotesStreamRequest) (*GetSumVotesStreamResponse, error) {
	return &GetSumVotesStreamResponse{}, nil
}
