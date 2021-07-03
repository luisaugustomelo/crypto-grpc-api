package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

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

// var dbClient *mongo.Client
var mongoCtx context.Context
var db *mongo.Collection

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

	log.Print("Connected to mongodb successfully")
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

	crypto.Id = result.InsertedID.(primitive.ObjectID).Hex()

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
		Name:        name,
		Initials:    initials,
		Description: description,
	}

	result := db.FindOneAndUpdate(mongoCtx, bson.M{"_id": cryptoId}, bson.M{"$set": data}, options.FindOneAndUpdate().SetReturnDocument(1))

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

	if err != nil && result.DeletedCount == 0 {
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
		return nil, status.Error(codes.NotFound, "Cannot be delete a crypto with this Object Id")
	}

	return &system.ReadCryptocurrencyResponse{Crypto: &system.Cryptocurrency{
		Id:          crypto.Id.Hex(),
		Name:        crypto.Name,
		Initials:    crypto.Initials,
		Upvote:      crypto.Upvote,
		Downvote:    crypto.Downvote,
		Description: crypto.Description,
	}}, nil
}

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

func broadcast(msg models.Cryptocurrency) {
	for _, channel := range allRegisteredClients {
		select {
		case channel <- msg:
		default:
		}
	}

}

func (s *Server) UpVoteCriptocurrency(ctx context.Context, request *system.UpVoteCryptocurrencyRequest) (*system.UpVoteCryptocurrencyResponse, error) {
	cryptoID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	filter := bson.M{"_id": cryptoID}

	result := db.FindOneAndUpdate(
		mongoCtx,
		filter,
		bson.M{"$inc": bson.M{"Upvote": 1}},
		options.FindOneAndUpdate().SetReturnDocument(1),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Couldn`t find Cryptocurrency with Object Id")
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
			Description: crypto.Description,
			Downvote:    crypto.Downvote,
			Upvote:      crypto.Upvote,
		},
	}

	return response, nil
}

func (s *Server) DownVoteCriptocurrency(ctx context.Context, request *system.DownVoteCryptocurrencyRequest) (*system.DownVoteCryptocurrencyResponse, error) {
	cryptoID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	filter := bson.M{"_id": cryptoID}

	result := db.FindOneAndUpdate(
		mongoCtx,
		filter,
		bson.M{"$inc": bson.M{"Downvote": 1}},
		options.FindOneAndUpdate().SetReturnDocument(1),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Couldn`t find Cryptocurrency with Object Id")
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
			Description: crypto.Description,
			Downvote:    int32(crypto.Downvote),
			Upvote:      crypto.Upvote,
		},
	}

	return response, nil
}

func (s *Server) GetSumVotes(ctx context.Context, request *system.GetSumVotesRequest) (*system.GetSumVotesResponse, error) {
	cryptoID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	result := db.FindOne(mongoCtx, bson.M{"_id": cryptoID})

	data := models.Cryptocurrency{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, "Couldn`t find Cryptocurrency with Object Id")
	}

	response := &system.GetSumVotesResponse{
		Votes: data.Upvote - data.Downvote,
	}
	return response, nil
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

func (s *Server) GetSumVotesByStream(request *system.GetSumVotesStreamRequest, stream system.UpVoteService_GetSumVotesByStreamServer) error {
	cryptoID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	result := db.FindOne(mongoCtx, bson.M{"_id": cryptoID})

	data := models.Cryptocurrency{}

	if err := result.Decode(&data); err != nil {
		return status.Errorf(codes.NotFound, "Couldn`t find Cryptocurrency with Object Id")
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
		if cryptoID == crypto.Id {
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
