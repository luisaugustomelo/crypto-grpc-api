package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	mongo "klever/grpc/databases"
	"klever/grpc/databases/config"

	crud "klever/grpc/services/cryptocurrencies"
	manage "klever/grpc/services/votes"

	"github.com/joho/godotenv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"klever/grpc/upvote/system"
)

var db mongo.CollectionHelper
var mongoCtx context.Context

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

func (s *Server) CreateCryptocurrency(ctx context.Context, request *system.CreateCryptocurrencyRequest) (*system.CreateCryptocurrencyResponse, error) {
	response, err := crud.CreateCryptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) UpdateCryptocurrency(ctx context.Context, request *system.UpdateCryptocurrencyRequest) (*system.UpdateCryptocurrencyResponse, error) {
	response, err := crud.UpdateCryptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) DeleteCryptocurrency(ctx context.Context, request *system.DeleteCryptocurrencyRequest) (*system.DeleteCryptocurrencyResponse, error) {
	response, err := crud.DeleteCryptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) ReadCryptocurrencyById(ctx context.Context, request *system.ReadCryptocurrencyRequest) (*system.ReadCryptocurrencyResponse, error) {
	response, err := crud.ReadCryptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) ListAllCriptocurrencies(request *system.ListAllCryptocurrenciesRequest, stream system.UpVoteService_ListAllCriptocurrenciesServer) error {
	err := manage.ListAllCriptocurrenciesService(request, stream, db, mongoCtx)

	return err
}

func (s *Server) UpVoteCriptocurrency(ctx context.Context, request *system.UpVoteCryptocurrencyRequest) (*system.UpVoteCryptocurrencyResponse, error) {
	response, err := manage.UpVoteCriptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) DownVoteCriptocurrency(ctx context.Context, request *system.DownVoteCryptocurrencyRequest) (*system.DownVoteCryptocurrencyResponse, error) {
	response, err := manage.DownVoteCriptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) GetSumVotes(ctx context.Context, request *system.GetSumVotesRequest) (*system.GetSumVotesResponse, error) {
	response, err := manage.GetSumVotesCryptocurrencyService(request, db, mongoCtx)

	return response, err
}

func (s *Server) GetSumVotesByStream(request *system.GetSumVotesStreamRequest, stream system.UpVoteService_GetSumVotesByStreamServer) error {
	err := manage.GetSumVotesByStreamService(request, stream, db, mongoCtx)

	return err
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
