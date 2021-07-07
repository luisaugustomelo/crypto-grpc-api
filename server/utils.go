package server

import (
	"fmt"
	mongo "klever/grpc/databases"
	"klever/grpc/databases/config"
	"klever/grpc/models"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

func loadEnvironmentVariable() {
	path, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(path, ".env"))

	if err != nil {
		log.Fatalf("Error %s loading .env file", err)
	}
}

func connectoToMongoDB() {

	conf := config.GetConfig()
	dbClient, err := mongo.NewClient(conf)

	if err != nil {
		log.Fatalf("Failed to create new database client: %s", err)
	}

	mongoCtx = context.Background()
	err = dbClient.Connect(mongoCtx)

	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}

	db = dbClient.Database(conf.DatabaseName).Collection(conf.Collection)

	log.Print("Connected to mongodb successfully")
}

func broadcast(msg models.Cryptocurrency) {
	for _, channel := range allRegisteredClients {
		select {
		case channel <- msg:
		default:
		}
	}
}

func disconectClient(channel chan models.Cryptocurrency) {
	removeClientMutex.Lock()
	defer removeClientMutex.Unlock()
	found := false
	i := 0

	for ; i < len(allRegisteredClients); i++ {
		if allRegisteredClients[i] == channel {
			found = true
			break
		}
	}
	if found {
		allRegisteredClients[i] = allRegisteredClients[len(allRegisteredClients)-1]
		allRegisteredClients = allRegisteredClients[:len(allRegisteredClients)-1]
	}
}

func connectionDb() {
	path, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(path, "..", ".env"))

	if err != nil {
		log.Fatalf("Error %s loading .env file", err)
	}

	conf := config.GetConfig()
	dbClient, err := mongo.NewClient(conf)

	if err != nil {
		log.Fatal(err)
	}

	mongoCtx = context.Background()
	err = dbClient.Connect(mongoCtx)

	if err != nil {
		log.Fatal(err)
	}

	db = dbClient.Database("klever").Collection("test")
	fmt.Println("Connected to MongoDB")
}

func clearDB() {
	db.Drop(mongoCtx)
}
