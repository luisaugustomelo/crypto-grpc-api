package manage

import (
	"context"
	"log"
	"sync"
	"time"

	// mongo "klever/grpc/databases"
	mongo "klever/grpc/databases"
	"klever/grpc/models"
	system "klever/grpc/upvote/system"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var allRegisteredClients []chan models.Cryptocurrency

var removeClientMutex sync.Mutex

func UpVoteCriptocurrencyService(request *system.UpVoteCryptocurrencyRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.UpVoteCryptocurrencyResponse, error) {
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

func DownVoteCriptocurrencyService(request *system.DownVoteCryptocurrencyRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.DownVoteCryptocurrencyResponse, error) {
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

func broadcast(msg models.Cryptocurrency) {
	for _, channel := range allRegisteredClients {
		select {
		case channel <- msg:
		default:
		}
	}
}

func GetSumVotesCryptocurrencyService(request *system.GetSumVotesRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.GetSumVotesResponse, error) {
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

func GetSumVotesByStreamService(request *system.GetSumVotesStreamRequest, stream system.UpVoteService_GetSumVotesByStreamServer, db mongo.CollectionHelper, mongoCtx context.Context) error {
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

func disconectClient(channel chan models.Cryptocurrency) {
	removeClientMutex.Lock()
	defer removeClientMutex.Unlock()
	found := false
	i := 0

	for ; i < len(allRegisteredClients); i++ {
		if allRegisteredClients[i] == channel {
			found = true
			break
		}
	}
	if found {
		allRegisteredClients[i] = allRegisteredClients[len(allRegisteredClients)-1]
		allRegisteredClients = allRegisteredClients[:len(allRegisteredClients)-1]
	}
}

func ListAllCriptocurrenciesService(request *system.ListAllCryptocurrenciesRequest, stream system.UpVoteService_ListAllCriptocurrenciesServer, db mongo.CollectionHelper, mongoCtx context.Context) error {
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
