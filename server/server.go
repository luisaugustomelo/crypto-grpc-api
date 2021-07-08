package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	mongo "klever/grpc/databases"
	"klever/grpc/databases/config"
	"klever/grpc/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"klever/grpc/proto/system"
)

var db mongo.CollectionHelper
var mongoCtx context.Context
var allRegisteredClients []chan models.Cryptocurrency

var removeClientMutex sync.Mutex

type Server struct {
	system.UnimplementedUpVoteServiceServer
}

func loadEnvironmentVariable() {
	path, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(path, ".env"))

	if err != nil {
		log.Fatalf("Error %s loading .env file", err)
	}
}

func connectoToMongoDB() {

	conf := config.GetConfig()
	dbClient, err := mongo.NewClient(conf)

	if err != nil {
		log.Fatalf("Failed to create new database client: %s", err)
	}

	mongoCtx = context.Background()
	err = dbClient.Connect(mongoCtx)

	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}

	db = dbClient.Database(conf.DatabaseName).Collection(conf.Collection)

	log.Print("Connected to mongodb successfully")
}

func (s *Server) HealthCheck(ctx context.Context, message *system.Message) (*system.Message, error) {
	log.Printf("Received message of health check from client: %v", message.Body)

	return &system.Message{Body: "****** The server it's OK! ******"}, nil
}

func (s *Server) CleanCollection(ctx context.Context, message *system.CleanCollectionRequest) (*system.CleanCollectionResponse, error) {
	deletedCount, err := db.Drop(ctx)

	return &system.CleanCollectionResponse{DeletedCount: deletedCount}, err
}

func main() {
	loadEnvironmentVariable()

	connectoToMongoDB()

	port := os.Getenv("KLEVER_APPLICATION_PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %s", port, err)
	} else {
		log.Printf("Application listening on port %s", port)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	system.RegisterUpVoteServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}

// CRUD

func (s *Server) CreateCryptocurrency(ctx context.Context, request *system.CreateCryptocurrencyRequest) (*system.CreateCryptocurrencyResponse, error) {
	crypto := request.GetCrypto()

	name := crypto.GetName()
	description := crypto.GetDescription()
	initials := crypto.GetInitials()

	if name == "" || description == "" || initials == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty fields")
	}

	data := models.Cryptocurrency{
		Id:          primitive.NewObjectID(),
		Name:        name,
		Initials:    initials,
		Upvote:      0,
		Downvote:    0,
		Description: description,
	}

	findResult := db.FindOne(mongoCtx, bson.M{"name": name})

	isRegisteredCrypto := models.Cryptocurrency{}

	if err := findResult.Decode(&isRegisteredCrypto); err == nil {
		return nil, status.Error(codes.AlreadyExists, "Cryptocurrency already exists")
	}
	result, err := db.InsertOne(mongoCtx, data)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	crypto.Id = result.(primitive.ObjectID).Hex()

	response := &system.CreateCryptocurrencyResponse{Crypto: crypto}

	return response, nil
}

func (s *Server) UpdateCryptocurrency(ctx context.Context, request *system.UpdateCryptocurrencyRequest) (*system.UpdateCryptocurrencyResponse, error) {
	crypto := request.GetCrypto()

	//Check if id is valid
	cryptoId, err := primitive.ObjectIDFromHex(crypto.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	name := crypto.GetName()
	description := crypto.GetDescription()
	initials := crypto.GetInitials()

	if name == "" || description == "" || initials == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty fields")
	}

	data := models.Cryptocurrency{
		Id:          cryptoId,
		Name:        name,
		Initials:    initials,
		Description: description,
	}

	result := db.FindOneAndUpdate(mongoCtx, bson.M{"_id": cryptoId}, bson.M{"$set": data})

	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.NotFound, "Cannot be find a crypto with this Object Id")
	}

	isUpdatedCrypto := models.Cryptocurrency{}

	err = result.Decode(&isUpdatedCrypto)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
	}

	response := &system.UpdateCryptocurrencyResponse{Crypto: crypto}

	return response, nil
}

func (s *Server) DeleteCryptocurrency(ctx context.Context, request *system.DeleteCryptocurrencyRequest) (*system.DeleteCryptocurrencyResponse, error) {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := db.DeleteOne(mongoCtx, bson.M{"_id": cryptoId})

	if err != nil && result == 0 {
		return nil, status.Errorf(codes.NotFound, "Cannot be delete a crypto with this Object Id")
	}

	response := &system.DeleteCryptocurrencyResponse{Status: true}

	return response, nil
}

func (s *Server) ReadCryptocurrencyById(ctx context.Context, request *system.ReadCryptocurrencyRequest) (*system.ReadCryptocurrencyResponse, error) {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result := db.FindOne(mongoCtx, bson.M{"_id": cryptoId})

	crypto := models.Cryptocurrency{}

	if err := result.Decode(&crypto); err == nil {
		status.Error(codes.NotFound, "Cannot be find a crypto with this Object Id")
	}

	response := &system.ReadCryptocurrencyResponse{
		Crypto: &system.Cryptocurrency{
			Id:          crypto.Id.Hex(),
			Name:        crypto.Name,
			Initials:    crypto.Initials,
			Upvote:      crypto.Upvote,
			Downvote:    crypto.Downvote,
			Description: crypto.Description,
		}}

	return response, nil
}

//Votes

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

//Utils

func broadcast(msg models.Cryptocurrency) {
	for _, channel := range allRegisteredClients {
		select {
		case channel <- msg:
		default:
		}
	}
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

func connectionDb() {
	path, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(path, "..", ".env"))

	if err != nil {
		log.Fatalf("Error %s loading .env file", err)
	}

	conf := config.GetConfig()
	dbClient, err := mongo.NewClient(conf)

	if err != nil {
		log.Fatal(err)
	}

	mongoCtx = context.Background()
	err = dbClient.Connect(mongoCtx)

	if err != nil {
		log.Fatal(err)
	}

	db = dbClient.Database("klever").Collection("test")

	fmt.Println("Connected to MongoDB")
}
