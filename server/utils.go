package main

import (
	mongo "klever/grpc/databases"
	"klever/grpc/databases/config"
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
