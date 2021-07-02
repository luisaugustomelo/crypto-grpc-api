package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"klever/grpc/models"
	"klever/grpc/upvote/system"
)

var dbClient *mongo.Client
var mongoCtx context.Context
var db *mongo.Collection

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
	db_port := os.Getenv("KLEVER_MONGODB_PORT")

	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:" + db_port))

	if err != nil {
		log.Fatalf("Problem with mongodb %s", err)
	}

	mongoCtx = context.Background()
	err = dbClient.Connect(mongoCtx)

	if err != nil {
		log.Fatalf("Error to connect from mongodb %s", err)
	}

	db = dbClient.Database("klever").Collection("cryptocurrencies")

	log.Print("Connected to mongodb successly")
}

func (s *Server) PingPong(ctx context.Context, message *system.Message) (*system.Message, error) {
	log.Printf("Received message body from client: %v", message.Body)

	return &system.Message{Body: "Hello From the Server!"}, nil
}

func (s *Server) CreateCryptocurrency(ctx context.Context, request *system.CreateCryptocurrencyRequest) (*system.CreateCryptocurrencyResponse, error) {
	crypto := request.GetCrypto()

	name := crypto.GetName()
	description := crypto.GetDescription()
	initials := crypto.GetInitials()

	if name == "" || description == "" {
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

	cryptoDUp := models.Cryptocurrency{}

	if err := findResult.Decode(&cryptoDUp); err == nil {
		return nil, status.Error(codes.AlreadyExists, "Cryptocurrency already exists")
	}

	insertResult, err := db.InsertOne(mongoCtx, data)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	crypto.Id = insertResult.InsertedID.(primitive.ObjectID).Hex()

	response := &system.CreateCryptocurrencyResponse{Crypto: crypto}

	return response, nil
}

func (s *Server) UpdateCryptocurrency(ctx context.Context, request *system.UpdateCryptocurrencyRequest) (*system.UpdateCryptocurrencyResponse, error) {
	return &system.UpdateCryptocurrencyResponse{}, nil
}

func (s *Server) DeleteCryptocurrency(ctx context.Context, request *system.DeleteCryptocurrencyRequest) (*system.DeleteCryptocurrencyResponse, error) {
	return &system.DeleteCryptocurrencyResponse{}, nil
}

func (s *Server) ReadCryptocurrencyById(ctx context.Context, request *system.ReadCryptocurrencyRequest) (*system.ReadCryptocurrencyResponse, error) {
	return &system.ReadCryptocurrencyResponse{}, nil
}

func (s *Server) ListAllCriptocurrencies(ctx context.Context, request *system.ListAllCryptocurrenciesRequest) (*system.ListAllCryptocurrenciesResponse, error) {
	return &system.ListAllCryptocurrenciesResponse{}, nil
}

func (s *Server) UpVoteCriptocurrency(ctx context.Context, request *system.UpVoteCryptocurrencyRequest) (*system.UpVoteCryptocurrencyResponse, error) {
	return &system.UpVoteCryptocurrencyResponse{}, nil
}

func (s *Server) DownVoteCriptocurrency(ctx context.Context, request *system.DownVoteCryptocurrencyRequest) (*system.DownVoteCryptocurrencyResponse, error) {
	return &system.DownVoteCryptocurrencyResponse{}, nil
}

func (s *Server) GetSumVotes(ctx context.Context, request *system.GetSumVotesRequest) (*system.GetSumVotesResponse, error) {
	return &system.GetSumVotesResponse{}, nil
}

func (s *Server) GetSumVotesByStream(request *system.GetSumVotesStreamRequest, stream system.UpVoteService_GetSumVotesByStreamServer) error {
	return nil
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
