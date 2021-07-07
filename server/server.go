package server

import (
	"log"
	"net"
	"os"

	mongo "klever/grpc/databases"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"klever/grpc/upvote/system"
)

var db mongo.CollectionHelper
var mongoCtx context.Context

type Server struct {
	system.UnimplementedUpVoteServiceServer
}

func (s *Server) HealthCheck(ctx context.Context, message *system.Message) (*system.Message, error) {
	log.Printf("Received message of health check from client: %v", message.Body)

	return &system.Message{Body: "****** The server it's OK! ******"}, nil
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
