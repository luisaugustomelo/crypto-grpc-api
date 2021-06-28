package main

import (
	"klever/grpc/upvote/klever"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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

func main() {
	port := goDotEnvVariable("KLEVER_APPLICATION_PORT")

	lis, err := net.Listen("tcp", port)
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
