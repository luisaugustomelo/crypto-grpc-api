package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	system "klever/grpc/upvote/system"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	message := system.Message{
		Body: "Hello from the client!",
	}

	response, err := c.PingPong(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling UpVote: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)

}
