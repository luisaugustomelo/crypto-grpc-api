package client

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
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

	//create
	createCrypto := system.CreateCryptocurrencyRequest{Crypto: &crypto}

	createResponse, err := c.CreateCryptocurrency(context.Background(), &createCrypto)
	if err != nil {
		log.Fatalf("Error when create a cryptocurrency: %s", err)
	}

	log.Printf("Crypto created as success: %s", createResponse.Crypto)

	//update
	crypto.Description = "...."
	crypto.Id = createResponse.Crypto.Id

	updateResponse, err := c.UpdateCryptocurrency(context.Background(), &system.UpdateCryptocurrencyRequest{Crypto: &crypto})

	if err != nil {
		log.Fatalf("Error when update the cryptocurrency: %s", err)
	}

	log.Printf("Crypto updated as success: %s", updateResponse.Crypto)

	//delete

	deleteResponse, err := c.DeleteCryptocurrency(context.Background(), &system.DeleteCryptocurrencyRequest{Id: crypto.Id})

	if err != nil {
		log.Fatalf("Error when update the cryptocurrency: %s", err)
	}

	log.Printf("Crypto deleted as success: %s", deleteResponse)

	//create

	crypto = system.Cryptocurrency{
		Name:        "Tron",
		Initials:    "TRX",
		Description: "TRON was founded by Justin Sun, who now serves as CEO. Educated at Peking University and the University of Pennsylvania, he was recognized by Forbes Asia in its 30 Under 30 series for entrepreneurs. Born in 1990, he was also associated with Ripple in the past — serving as its chief representative in the Greater China area.",
		Upvote:      0,
		Downvote:    0,
	}

	createCrypto = system.CreateCryptocurrencyRequest{Crypto: &crypto}

	createResponse, err = c.CreateCryptocurrency(context.Background(), &createCrypto)
	if err != nil {
		log.Fatalf("Error when create a cryptocurrency: %s", err)
	}

	log.Printf("Crypto created as success: %s", createResponse.Crypto)

	//read

	log.Printf("Crypto to read %s", createResponse.Crypto.Id)

	readResponse, err := c.ReadCryptocurrencyById(context.Background(), &system.ReadCryptocurrencyRequest{Id: createResponse.Crypto.Id})

	if err != nil {
		log.Fatalf("Error when read the cryptocurrency: %s", err)
	}

	log.Printf("Crypto read as success: %s", readResponse)

	//read all

	log.Print("Reading all cryptocurrencies")

	request := &system.ListAllCryptocurrenciesRequest{}

	stream, err := c.ListAllCriptocurrencies(context.Background(), request)

	if err != nil {
		log.Fatalf("Error to read all cryptoes")
	}

	result := getAllCryptos(stream)

	log.Printf("%s", gin.H{
		"result": result,
	})

	//Upvote

	log.Print("Upvote Tron(TRX) crypto")

	upVoteResponse, err := c.UpVoteCriptocurrency(context.Background(), &system.UpVoteCryptocurrencyRequest{Id: createResponse.Crypto.Id})

	if err != nil {
		log.Fatalf("Error when read the cryptocurrency: %s", err)
	}

	log.Printf("Crypto voted as success: %s", upVoteResponse)

	upVoteResponse, err = c.UpVoteCriptocurrency(context.Background(), &system.UpVoteCryptocurrencyRequest{Id: createResponse.Crypto.Id})
	if err != nil {
		log.Fatalf("Error when read the cryptocurrency: %s", err)
	}

	log.Printf("Crypto voted as success: %s", upVoteResponse)

	upVoteResponse, err = c.UpVoteCriptocurrency(context.Background(), &system.UpVoteCryptocurrencyRequest{Id: createResponse.Crypto.Id})
	if err != nil {
		log.Fatalf("Error when read the cryptocurrency: %s", err)
	}

	log.Printf("Crypto voted as success: %s", upVoteResponse)

	//Downvote

	log.Print("Downvote Tron(TRX) crypto")

	downVoteResponse, err := c.DownVoteCriptocurrency(context.Background(), &system.DownVoteCryptocurrencyRequest{Id: createResponse.Crypto.Id})

	if err != nil {
		log.Fatalf("Error when read the cryptocurrency: %s", err)
	}

	log.Printf("Crypto voted as success: %s", downVoteResponse)

	//GetSumVotes

	getSumVotes, err := c.GetSumVotes(context.Background(), &system.GetSumVotesRequest{Id: createResponse.Crypto.Id})

	if err != nil {
		log.Fatalf("Error when read the cryptocurrency: %s", err)
	}

	log.Printf("Crypto sum votes: %s", getSumVotes)

	//read sum by stream

	// log.Print("Reading all cryptocurrencies")

	// request := &system.GetSumVotesStreamRequest{Id: createResponse.Crypto.Id}

	// stream, err := c.GetSumVotesByStream(context.Background(), request)

	// if err != nil {
	// 	log.Fatalf("Error to read all cryptoes")
	// }

	// result := getAllCryptos(stream)

	// log.Printf("%s", gin.H{
	// 	"result": result,
	// })

	//health check
	response, err := c.HealthCheck(context.Background(), &message)

	if err != nil {
		log.Fatalf("Error when calling UpVote: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)

}

func getAllCryptos(stream system.UpVoteService_ListAllCriptocurrenciesClient) []*system.ListAllCryptocurrenciesResponse {
	var processing []*system.ListAllCryptocurrenciesResponse

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print("Couldn`t get cryptocurrencies")
		}
		processing = append(processing, data)
	}

	return processing
}
