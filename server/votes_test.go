package main

import (
	"io"
	"klever/grpc/upvote/system"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// func TestUpvoteCrypto(t *testing.T) {
// 	connectionDb()
// 	var conn *grpc.ClientConn

// 	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Could not connect: %s", err)
// 	}

// 	defer conn.Close()

// 	c := system.NewUpVoteServiceClient(conn)

// 	defer clearDB()

// 	// Test request with empty ID
// 	emptyIDRequest := &system.UpVoteCryptocurrencyRequest{
// 		Id: "",
// 	}
// 	_, err = c.UpVoteCriptocurrency(context.Background(), emptyIDRequest)

// 	require.NotNil(t, err)

// 	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

// 	// Test request with valid ID but not found on DB
// 	NotFoundIDRequest := &system.UpVoteCryptocurrencyRequest{
// 		Id: primitive.NewObjectID().Hex(),
// 	}
// 	_, err = c.UpVoteCriptocurrency(context.Background(), NotFoundIDRequest)

// 	require.Nil(t, err)

// 	// Test with valid request
// 	createRequest := &system.CreateCryptocurrencyRequest{
// 		Crypto: &system.Cryptocurrency{
// 			Name:        "Bitcoin",
// 			Initials:    "BTC",
// 			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
// 		},
// 	}

// 	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), createRequest)

// 	require.Nil(t, err)

// 	validRequest := system.UpVoteCryptocurrencyRequest{
// 		Id: cryptoResponse.GetCrypto().GetId(),
// 	}
// 	response, err := c.UpVoteCriptocurrency(context.Background(), &validRequest)

// 	require.Nil(t, err)

// 	assert.Equal(t, cryptoResponse.GetCrypto().GetId(), response.GetCrypto().GetId())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetInitials(), response.GetCrypto().GetInitials())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetName(), response.GetCrypto().GetName())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetDescription(), response.GetCrypto().GetDescription())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetUpvote(), response.GetCrypto().GetUpvote())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetDownvote(), response.GetCrypto().GetDownvote())

// }

// func TestDownvoteCrypto(t *testing.T) {
// 	connectionDb()
// 	var conn *grpc.ClientConn

// 	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Could not connect: %s", err)
// 	}

// 	defer conn.Close()

// 	c := system.NewUpVoteServiceClient(conn)

// 	defer clearDB()

// 	// Test request with empty ID
// 	emptyIDRequest := &system.DownVoteCryptocurrencyRequest{
// 		Id: "",
// 	}
// 	_, err = c.DownVoteCriptocurrency(context.Background(), emptyIDRequest)

// 	require.NotNil(t, err)

// 	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

// 	// Test request with valid ID but not found on DB
// 	NotFoundIDRequest := &system.DownVoteCryptocurrencyRequest{
// 		Id: primitive.NewObjectID().Hex(),
// 	}
// 	_, err = c.DownVoteCriptocurrency(context.Background(), NotFoundIDRequest)

// 	require.Nil(t, err)

// 	// Test with valid request
// 	createRequest := system.CreateCryptocurrencyRequest{
// 		Crypto: &system.Cryptocurrency{
// 			Name:        "Bitcoin",
// 			Initials:    "BTC",
// 			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
// 		},
// 	}

// 	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), &createRequest)

// 	require.Nil(t, err)

// 	validRequest := &system.DownVoteCryptocurrencyRequest{
// 		Id: cryptoResponse.GetCrypto().GetId(),
// 	}
// 	response, err := c.DownVoteCriptocurrency(context.Background(), validRequest)

// 	require.Nil(t, err)

// 	assert.Equal(t, cryptoResponse.GetCrypto().GetId(), response.GetCrypto().GetId())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetInitials(), response.GetCrypto().GetInitials())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetName(), response.GetCrypto().GetName())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetDescription(), response.GetCrypto().GetDescription())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetUpvote(), response.GetCrypto().GetUpvote())
// 	assert.Equal(t, cryptoResponse.GetCrypto().GetDownvote(), response.GetCrypto().GetDownvote())

// }

// func TestGetVotesSum(t *testing.T) {
// 	connectionDb()
// 	var conn *grpc.ClientConn

// 	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Could not connect: %s", err)
// 	}

// 	defer conn.Close()

// 	c := system.NewUpVoteServiceClient(conn)

// 	defer clearDB()

// 	// Test request with empty ID
// 	emptyIDRequest := &system.GetSumVotesRequest{
// 		Id: "",
// 	}
// 	_, err = c.GetSumVotes(context.Background(), emptyIDRequest)

// 	require.NotNil(t, err)

// 	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

// 	// Test request with valid ID but not found on DB
// 	NotFoundIDRequest := &system.GetSumVotesRequest{
// 		Id: primitive.NewObjectID().Hex(),
// 	}
// 	_, err = c.GetSumVotes(context.Background(), NotFoundIDRequest)

// 	require.NotNil(t, err)

// 	assert.Equal(t, "rpc error: code = NotFound desc = Cannot find Cryptocurrency with Object Id", err.Error())

// 	// Test with valid request
// 	createRequest := system.CreateCryptocurrencyRequest{
// 		Crypto: &system.Cryptocurrency{
// 			Name:        "Bitcoin",
// 			Initials:    "BTC",
// 			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
// 		},
// 	}

// 	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), &createRequest)

// 	require.Nil(t, err)

// 	var downvoteResponse *system.DownVoteCryptocurrencyResponse

// 	for i := 0; i < 5; i++ {
// 		validRequest := &system.UpVoteCryptocurrencyRequest{
// 			Id: cryptoResponse.GetCrypto().GetId(),
// 		}
// 		_, err = c.UpVoteCriptocurrency(context.Background(), validRequest)

// 		require.Nil(t, err)
// 	}

// 	for i := 0; i < 2; i++ {
// 		validRequest := &system.DownVoteCryptocurrencyRequest{
// 			Id: cryptoResponse.GetCrypto().GetId(),
// 		}
// 		downvoteResponse, err = c.DownVoteCriptocurrency(context.Background(), validRequest)

// 		require.Nil(t, err)
// 	}

// 	validRequest := &system.GetSumVotesRequest{
// 		Id: cryptoResponse.GetCrypto().GetId(),
// 	}
// 	response, err := c.GetSumVotes(context.Background(), validRequest)

// 	require.Nil(t, err)

// 	assert.Equal(t, downvoteResponse.GetCrypto().GetUpvote()-downvoteResponse.GetCrypto().GetDownvote(), response.GetVotes()+1)

// }

func TestGetVotesSumStream(t *testing.T) {
	connectionDb()
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	// defer clearDB()

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Test invalid request
	client := system.NewUpVoteServiceClient(conn)
	request := &system.GetSumVotesStreamRequest{
		Id: "",
	}

	stream, err := client.GetSumVotesByStream(context.Background(), request)

	var result []*system.GetSumVotesStreamResponse
	var resp *system.GetSumVotesStreamResponse

	done := make(chan bool)

	go func() {
		for {
			resp, err = stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				done <- true
				return
			}

			result = append(result, resp)

		}
	}()

	<-done
	require.NotNil(t, err)
	require.Nil(t, result)

	// Test valid live update
	createRequest := &system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), createRequest)

	require.Nil(t, err)

	request = &system.GetSumVotesStreamRequest{
		Id: cryptoResponse.GetCrypto().GetId(),
	}

	stream, err = client.GetSumVotesByStream(context.Background(), request)

	go func() {
		for {
			resp, err = stream.Recv()
			if err == io.EOF || err != nil {
				break
			}
		}
	}()

	for i := 0; i < 4; i++ {
		upvoteRequest := &system.UpVoteCryptocurrencyRequest{
			Id: cryptoResponse.GetCrypto().GetId(),
		}
		_, err = c.UpVoteCriptocurrency(context.Background(), upvoteRequest)

		require.Nil(t, err)

	}

	for i := 0; i < 2; i++ {
		downvoteRequest := &system.DownVoteCryptocurrencyRequest{
			Id: cryptoResponse.GetCrypto().GetId(),
		}
		_, err = c.DownVoteCriptocurrency(context.Background(), downvoteRequest)

		require.Nil(t, err)

	}
	// Sleep to give stream time to receive before finish Test
	time.Sleep(200 * time.Microsecond)

	require.Nil(t, err)

	require.Equal(t, int32(2), resp.GetVotes()-1)

	clearDB()

}
