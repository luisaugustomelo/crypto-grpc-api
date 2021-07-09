package main_test

import (
	"io"
	"klever/grpc/proto/system"
	"log"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func TestCreateCryptocurrency(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	//Empty Fields
	crypto := system.Cryptocurrency{
		Name:        "",
		Initials:    "",
		Description: "",
	}

	createCrypto := system.CreateCryptocurrencyRequest{Crypto: &crypto}

	response, err := c.CreateCryptocurrency(context.Background(), &createCrypto)

	require.NotNil(t, err)
	require.Nil(t, response)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = Empty fields", err.Error())

	crypto = system.Cryptocurrency{
		Name:        "Bitcoin",
		Initials:    "BTC",
		Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
	}

	// Valid
	validRequest := system.CreateCryptocurrencyRequest{Crypto: &crypto}

	response, err = c.CreateCryptocurrency(context.Background(), &validRequest)

	require.Nil(t, err)

	assert.Equal(t, "Bitcoin", response.GetCrypto().GetName())
	assert.Equal(t, "BTC", response.GetCrypto().GetInitials())
	assert.Equal(t, "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)", response.GetCrypto().GetDescription())
	assert.Equal(t, int32(0), response.GetCrypto().GetUpvote())
	assert.Equal(t, int32(0), response.GetCrypto().GetDownvote())

	// Already exists
	response, err = c.CreateCryptocurrency(context.Background(), &validRequest)

	require.NotNil(t, err)
	require.Nil(t, response)
	assert.Equal(t, "rpc error: code = AlreadyExists desc = Cryptocurrency already exists", err.Error())
}

func TestUpdateCryptocurrency(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Empty ID
	emptyRequest := &system.UpdateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Id: "",
		},
	}
	_, err = c.UpdateCryptocurrency(context.Background(), emptyRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

	// Empty fields
	emptyFieldsRequest := &system.UpdateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Id:          primitive.NewObjectID().Hex(),
			Name:        "",
			Description: "",
		},
	}
	_, err = c.UpdateCryptocurrency(context.Background(), emptyFieldsRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = Empty fields", err.Error())

	// Valid but ID Not found on the Mongo
	notFoundRequest := &system.UpdateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Id:          primitive.NewObjectID().Hex(),
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}
	_, err = c.UpdateCryptocurrency(context.Background(), notFoundRequest)

	require.Nil(t, err)

	// Valid
	createRequest := &system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), createRequest)

	require.Nil(t, err)

	validRequest := system.UpdateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Id:          cryptoResponse.GetCrypto().GetId(),
			Name:        "Ethereum",
			Initials:    "BTC",
			Description: "Ethereum is a decentralized open-source blockchain system that features its own cryptocurrency (...)",
		},
	}

	response, err := c.UpdateCryptocurrency(context.Background(), &validRequest)

	require.Nil(t, err)

	//Check fields response
	assert.Equal(t, validRequest.GetCrypto().GetId(), response.GetCrypto().GetId())
	assert.Equal(t, validRequest.GetCrypto().GetInitials(), response.GetCrypto().GetInitials())
	assert.Equal(t, validRequest.GetCrypto().GetName(), response.GetCrypto().GetName())
	assert.Equal(t, validRequest.GetCrypto().GetDescription(), response.GetCrypto().GetDescription())
	assert.Equal(t, validRequest.GetCrypto().GetUpvote(), response.GetCrypto().GetUpvote())
	assert.Equal(t, validRequest.GetCrypto().GetDownvote(), response.GetCrypto().GetDownvote())
}

func TestDeleteCryptocurrency(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Empty
	emptyRequest := &system.DeleteCryptocurrencyRequest{
		Id: "",
	}
	_, err = c.DeleteCryptocurrency(context.Background(), emptyRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

	// Valid but ID Not found on the Mongo
	notFoundRequest := system.DeleteCryptocurrencyRequest{
		Id: primitive.NewObjectID().Hex(),
	}

	_, err = c.DeleteCryptocurrency(context.Background(), &notFoundRequest)

	require.Nil(t, err)

	// Valid
	createRequest := &system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), createRequest)

	require.Nil(t, err)

	validRequest := &system.DeleteCryptocurrencyRequest{
		Id: cryptoResponse.Crypto.GetId(),
	}

	response, err := c.DeleteCryptocurrency(context.Background(), validRequest)

	require.Nil(t, err)

	assert.Equal(t, true, response.GetStatus())
}

func TestReadCryptocurrencyById(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Empty
	emptyRequest := &system.ReadCryptocurrencyRequest{
		Id: "",
	}
	response, err := c.ReadCryptocurrencyById(context.Background(), emptyRequest)

	require.NotNil(t, err)
	require.Nil(t, response)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

	// Valid but ID Not found on the Mongo
	notFoundRequest := &system.ReadCryptocurrencyRequest{
		Id: primitive.NewObjectID().Hex(),
	}

	response, err = c.ReadCryptocurrencyById(context.Background(), notFoundRequest)

	require.Nil(t, err)
	require.NotNil(t, response)

	// Valid
	createRequest := system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), &createRequest)

	require.Nil(t, err)

	validRequest := system.ReadCryptocurrencyRequest{
		Id: cryptoResponse.Crypto.GetId(),
	}

	response, err = c.ReadCryptocurrencyById(context.Background(), &validRequest)

	require.NotNil(t, response)
	require.Nil(t, err)

	assert.Equal(t, cryptoResponse.GetCrypto().GetId(), response.GetCrypto().GetId())
	assert.Equal(t, cryptoResponse.GetCrypto().GetName(), response.GetCrypto().GetName())
	assert.Equal(t, cryptoResponse.GetCrypto().GetInitials(), response.GetCrypto().GetInitials())
	assert.Equal(t, cryptoResponse.GetCrypto().GetDescription(), response.GetCrypto().GetDescription())
	assert.Equal(t, cryptoResponse.GetCrypto().GetUpvote(), response.GetCrypto().GetUpvote())
	assert.Equal(t, cryptoResponse.GetCrypto().GetDownvote(), response.GetCrypto().GetDownvote())
}

func TestListAllCriptocurrencies(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Not found on the Mongo
	var allCreatedCrypto []*system.CreateCryptocurrencyResponse

	client := system.NewUpVoteServiceClient(conn)
	request := &system.ListAllCryptocurrenciesRequest{}

	stream, err := client.ListAllCriptocurrencies(context.Background(), request)

	var result []*system.ListAllCryptocurrenciesResponse
	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			require.Nil(t, err)

			result = append(result, resp)

		}
	}()

	<-done

	require.Nil(t, err)
	require.Nil(t, result)

	a := system.Cryptocurrency{
		Name:        "Iota",
		Initials:    "MIOTA",
		Description: "IOTA is a distributed ledger with one big difference: it isn’t actually a blockchain.",
	}

	b := system.Cryptocurrency{
		Name:        "Klever",
		Initials:    "KLV",
		Description: "Klever (KLV) is a crypto wallet ecosystem serving above 2 million users globally with Klever App (...)",
	}

	d := system.Cryptocurrency{
		Name:        "Tether",
		Initials:    "USDT",
		Description: "USDT is a stablecoin (stable-value cryptocurrency) that mirrors the price of the U.S. dollar (...)",
	}

	//mock
	cryptoTest := []*system.Cryptocurrency{&a, &b, &d}

	for _, crypto := range cryptoTest {
		createRequest := &system.CreateCryptocurrencyRequest{Crypto: crypto}
		response, err := c.CreateCryptocurrency(context.Background(), createRequest)
		require.Nil(t, err)
		allCreatedCrypto = append(allCreatedCrypto, response)
		require.NotNil(t, allCreatedCrypto)
	}

	stream, err = client.ListAllCriptocurrencies(context.Background(), request)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.Nil(t, err)

		result = append(result, resp)

	}

	require.Nil(t, err)
	for i, crypto := range cryptoTest {
		assert.Equal(t, crypto.GetName(), result[i].GetCrypto().GetName())
		assert.Equal(t, crypto.GetInitials(), result[i].GetCrypto().GetInitials())
		assert.Equal(t, crypto.GetDescription(), result[i].GetCrypto().GetDescription())
		assert.Equal(t, crypto.GetUpvote(), result[i].GetCrypto().GetUpvote())
		assert.Equal(t, crypto.GetDownvote(), result[i].GetCrypto().GetDownvote())
	}
}

func TestGetSumVotesByStream(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Invalid
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

	// Valid
	createRequest := system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), &createRequest)

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

	time.Sleep(200 * time.Microsecond)

	require.Nil(t, err)

	require.Equal(t, int32(2), resp.GetVotes()-1)
}

func TestUpVoteCriptocurrency(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Empty
	emptyRequest := &system.UpVoteCryptocurrencyRequest{
		Id: "",
	}
	_, err = c.UpVoteCriptocurrency(context.Background(), emptyRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

	// Valid but ID Not found on the Mongo
	notFoundRequest := &system.UpVoteCryptocurrencyRequest{
		Id: primitive.NewObjectID().Hex(),
	}
	_, err = c.UpVoteCriptocurrency(context.Background(), notFoundRequest)

	require.Nil(t, err)

	// Valid
	createRequest := &system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), createRequest)

	require.Nil(t, err)

	validRequest := system.UpVoteCryptocurrencyRequest{
		Id: cryptoResponse.GetCrypto().GetId(),
	}
	response, err := c.UpVoteCriptocurrency(context.Background(), &validRequest)

	require.Nil(t, err)

	assert.Equal(t, cryptoResponse.GetCrypto().GetId(), response.GetCrypto().GetId())
	assert.Equal(t, cryptoResponse.GetCrypto().GetInitials(), response.GetCrypto().GetInitials())
	assert.Equal(t, cryptoResponse.GetCrypto().GetName(), response.GetCrypto().GetName())
	assert.Equal(t, cryptoResponse.GetCrypto().GetDescription(), response.GetCrypto().GetDescription())
	assert.Equal(t, cryptoResponse.GetCrypto().GetUpvote(), response.GetCrypto().GetUpvote())
	assert.Equal(t, cryptoResponse.GetCrypto().GetDownvote(), response.GetCrypto().GetDownvote())
}

func TestDownVoteCriptocurrency(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Empty
	emptyRequest := &system.DownVoteCryptocurrencyRequest{
		Id: "",
	}
	_, err = c.DownVoteCriptocurrency(context.Background(), emptyRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

	// Valid but ID Not found on the Mongo
	notFoundRequest := &system.DownVoteCryptocurrencyRequest{
		Id: primitive.NewObjectID().Hex(),
	}
	_, err = c.DownVoteCriptocurrency(context.Background(), notFoundRequest)

	require.Nil(t, err)

	// Valid
	createRequest := system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), &createRequest)

	require.Nil(t, err)

	validRequest := &system.DownVoteCryptocurrencyRequest{
		Id: cryptoResponse.GetCrypto().GetId(),
	}
	response, err := c.DownVoteCriptocurrency(context.Background(), validRequest)

	require.Nil(t, err)

	assert.Equal(t, cryptoResponse.GetCrypto().GetId(), response.GetCrypto().GetId())
	assert.Equal(t, cryptoResponse.GetCrypto().GetInitials(), response.GetCrypto().GetInitials())
	assert.Equal(t, cryptoResponse.GetCrypto().GetName(), response.GetCrypto().GetName())
	assert.Equal(t, cryptoResponse.GetCrypto().GetDescription(), response.GetCrypto().GetDescription())
	assert.Equal(t, cryptoResponse.GetCrypto().GetUpvote(), response.GetCrypto().GetUpvote())
	assert.Equal(t, cryptoResponse.GetCrypto().GetDownvote(), response.GetCrypto().GetDownvote())
}

func TestGetSumVotes(t *testing.T) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	c := system.NewUpVoteServiceClient(conn)

	c.CleanCollection(context.Background(), &system.CleanCollectionRequest{})

	// Empty
	emptyRequest := &system.GetSumVotesRequest{
		Id: "",
	}
	_, err = c.GetSumVotes(context.Background(), emptyRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = the provided hex string is not a valid ObjectID", err.Error())

	// Valid but ID Not found on the Mongo
	notFoundRequest := system.GetSumVotesRequest{
		Id: primitive.NewObjectID().Hex(),
	}
	_, err = c.GetSumVotes(context.Background(), &notFoundRequest)

	require.NotNil(t, err)

	assert.Equal(t, "rpc error: code = NotFound desc = Cannot find Cryptocurrency with Object Id", err.Error())

	// Valid
	createRequest := system.CreateCryptocurrencyRequest{
		Crypto: &system.Cryptocurrency{
			Name:        "Bitcoin",
			Initials:    "BTC",
			Description: "Bitcoin is a decentralized cryptocurrency originally described in a 2008 (...)",
		},
	}

	cryptoResponse, err := c.CreateCryptocurrency(context.Background(), &createRequest)

	require.Nil(t, err)

	var downVoteResponse *system.DownVoteCryptocurrencyResponse

	for i := 0; i < 10; i++ {
		validRequest := system.UpVoteCryptocurrencyRequest{
			Id: cryptoResponse.GetCrypto().GetId(),
		}
		_, err = c.UpVoteCriptocurrency(context.Background(), &validRequest)

		require.Nil(t, err)
	}

	for i := 0; i < 8; i++ {
		validRequest := &system.DownVoteCryptocurrencyRequest{
			Id: cryptoResponse.GetCrypto().GetId(),
		}
		downVoteResponse, err = c.DownVoteCriptocurrency(context.Background(), validRequest)

		require.Nil(t, err)
	}

	validRequest := &system.GetSumVotesRequest{
		Id: cryptoResponse.GetCrypto().GetId(),
	}
	response, err := c.GetSumVotes(context.Background(), validRequest)

	require.Nil(t, err)

	assert.Equal(t, downVoteResponse.GetCrypto().GetUpvote()-downVoteResponse.GetCrypto().GetDownvote(), response.GetVotes()+1)

}
