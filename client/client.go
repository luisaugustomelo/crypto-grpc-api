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
		Body: "Hello, all right there?!",
	}

	crypto := system.Cryptocurrency{
		Name:        "Tron",
		Initials:    "TRX",
		Description: "TRON was founded by Justin Sun, who now serves as CEO. Educated at Peking University and the University of Pennsylvania, he was recognized by Forbes Asia in its 30 Under 30 series for entrepreneurs. Born in 1990, he was also associated with Ripple in the past — serving as its chief representative in the Greater China area.",
		Upvote:      0,
		Downvote:    0,
	}

	createCrypto := system.CreateCryptocurrencyRequest{Crypto: &crypto}

	createResponse, err := c.CreateCryptocurrency(context.Background(), &createCrypto)
	if err != nil {
		log.Fatalf("Error when create a cryptocurrency: %s", err)
	}

	log.Printf("Crypto created as success: %s", createResponse.Crypto)

	response, err := c.HealthCheck(context.Background(), &message)

	if err != nil {
		log.Fatalf("Error when calling UpVote: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)

}
