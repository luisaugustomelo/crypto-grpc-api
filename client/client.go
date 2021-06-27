package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"klever/grpc/upvote/klever"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := klever.NewUpVoteServiceClient(conn)

	message := klever.Message{
		Body: "Hello from the client!",
	}

	response, err := c.UpVote(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling UpVote: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)

}
