package server

import (
	mongo "klever/grpc/databases"
	"klever/grpc/models"
	"klever/grpc/upvote/system"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var allRegisteredClients []chan models.Cryptocurrency

var removeClientMutex sync.Mutex

func (s *Server) ListAllCriptocurrencies(request *system.ListAllCryptocurrenciesRequest, stream system.UpVoteService_ListAllCriptocurrenciesServer) error {
	data := &models.Cryptocurrency{}

	cryptos, err := db.Find(mongoCtx, bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, "Cannot find error: "+err.Error())
	}

	defer cryptos.Close(mongoCtx)

	for cryptos.Next(mongoCtx) {
		err := cryptos.Decode(data)
		if err != nil {
			return status.Errorf(codes.Unavailable, "Cannot be decode error: "+err.Error())
		}

		stream.Send(&system.ListAllCryptocurrenciesResponse{
			Crypto: &system.Cryptocurrency{
				Id:          data.Id.Hex(),
				Name:        data.Name,
				Initials:    data.Initials,
				Description: data.Description,
				Downvote:    data.Downvote,
				Upvote:      data.Upvote,
			},
		})
	}

	if err := cryptos.Err(); err != nil {
		return status.Errorf(codes.Internal, "Unkown mongoDB pointer error: "+err.Error())
	}

	return nil
}

func (s *Server) UpVoteCriptocurrency(ctx context.Context, request *system.UpVoteCryptocurrencyRequest) (*system.UpVoteCryptocurrencyResponse, error) {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	filter := bson.M{"_id": cryptoId}

	result := db.FindOneAndUpdate(
		mongoCtx,
		filter,
		bson.M{"$inc": bson.M{"upvote": 1}},
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "cannot find Cryptocurrency with Object Id")
		}
	}

	crypto := models.Cryptocurrency{}
	err = result.Decode(&crypto)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
	}

	broadcast(crypto)

	response := &system.UpVoteCryptocurrencyResponse{
		Crypto: &system.Cryptocurrency{
			Id:          crypto.Id.Hex(),
			Name:        crypto.Name,
			Initials:    crypto.Initials,
			Description: crypto.Description,
			Downvote:    crypto.Downvote,
			Upvote:      crypto.Upvote,
		},
	}

	return response, nil
}

func (s *Server) DownVoteCriptocurrency(ctx context.Context, request *system.DownVoteCryptocurrencyRequest) (*system.DownVoteCryptocurrencyResponse, error) {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	filter := bson.M{"_id": cryptoId}

	result := db.FindOneAndUpdate(
		mongoCtx,
		filter,
		bson.M{"$inc": bson.M{"downvote": 1}},
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Cannot find Cryptocurrency with Object Id")
		}
	}

	crypto := models.Cryptocurrency{}
	err = result.Decode(&crypto)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
	}

	broadcast(crypto)

	response := &system.DownVoteCryptocurrencyResponse{
		Crypto: &system.Cryptocurrency{
			Id:          crypto.Id.Hex(),
			Name:        crypto.Name,
			Initials:    crypto.Initials,
			Description: crypto.Description,
			Downvote:    int32(crypto.Downvote),
			Upvote:      crypto.Upvote,
		},
	}

	return response, nil
}

func (s *Server) GetSumVotes(ctx context.Context, request *system.GetSumVotesRequest) (*system.GetSumVotesResponse, error) {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	result := db.FindOne(mongoCtx, bson.M{"_id": cryptoId})

	data := models.Cryptocurrency{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, "Cannot find Cryptocurrency with Object Id")
	}

	response := &system.GetSumVotesResponse{
		Votes: data.Upvote - data.Downvote,
	}
	return response, nil
}

func (s *Server) GetSumVotesByStream(request *system.GetSumVotesStreamRequest, stream system.UpVoteService_GetSumVotesByStreamServer) error {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	result := db.FindOne(mongoCtx, bson.M{"_id": cryptoId})

	data := models.Cryptocurrency{}

	if err := result.Decode(&data); err != nil {
		return status.Errorf(codes.NotFound, "Cannot find Cryptocurrency with Object Id")
	}

	ch := make(chan models.Cryptocurrency)

	allRegisteredClients = append(allRegisteredClients, ch)

	streamCtx := stream.Context()
	go func() {
		for {
			if streamCtx.Err() == context.Canceled || streamCtx.Err() == context.DeadlineExceeded {

				disconectClient(ch)

				close(ch)
				log.Print("End stream")
				return
			}
			time.Sleep(time.Second)
		}

	}()

	for crypto := range ch {
		if cryptoId == crypto.Id {
			sum := crypto.Upvote - crypto.Downvote
			response := &system.GetSumVotesStreamResponse{
				Votes: sum,
			}
			err := stream.Send(response)
			if err != nil {

				disconectClient(ch)
				close(ch)
				log.Print("End stream")

				return nil
			}
		}
	}
	return nil
}
