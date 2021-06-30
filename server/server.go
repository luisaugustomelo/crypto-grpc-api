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

func goDotEnvVariable(key string) string {
	path, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(path, ".env"))

	if err != nil {
		log.Fatalf("Error %s loading .env file", err)
	}

	return os.Getenv(key)
}

func connectoToMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
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
	connectoToMongoDB()

	port := goDotEnvVariable("KLEVER_APPLICATION_PORT")

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
