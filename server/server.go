package main

import (
	"context"
	"klever/grpc/upvote/klever"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func loadEnvironmentVariable() {
	path, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(path, ".env"))

	if err != nil {
		log.Fatalf("Error %s loading .env file", err)
	}
}

func connectoToMongoDB() {
	db_port := os.Getenv("KLEVER_MONGODB_PORT")

	clientOptions := options.Client().ApplyURI("mongodb://mongodb:" + db_port)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatalf("Occured wrong behavior to connect on MongoDB: %s", err)
	}

	log.Printf("Connected to MongoDB on port: 27017!")
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

	s := klever.Server{}

	grpcServer := grpc.NewServer()

	klever.RegisterUpVoteServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}
