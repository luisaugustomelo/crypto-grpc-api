package main

import (
	// fake_mongo "klever/grpc/databases/test/fake_databases"

	system "klever/grpc/upvote/system"
	"log"
	"reflect"
	"testing"

	crud "klever/grpc/services/cryptocurrencies"
)

func TestCreateCryptocurrency(t *testing.T) {
	crypto := system.Cryptocurrency{
		Name:        "Tron",
		Initials:    "TRX",
		Description: "TRON was founded by Justin Sun, who now serves as CEO. Educated at Peking University and the University of Pennsylvania, he was recognized by Forbes Asia in its 30 Under 30 series for entrepreneurs. Born in 1990, he was also associated with Ripple in the past — serving as its chief representative in the Greater China area.",
		Upvote:      0,
		Downvote:    0,
	}

	//create
	createCrypto := system.CreateCryptocurrencyRequest{Crypto: &crypto}

	createResponse, err := crud.CreateCryptocurrencyService(&createCrypto, db, mongoCtx)

	if err != nil {
		log.Fatalf("Error when create a cryptocurrency: %s", err)
	}

	log.Printf("Crypto created as success: %s", createResponse.Crypto)
}
func TestUpVote(t *testing.T) {
	type Test struct {
		cryptocurrency string
		vote           int
	}

	s := []*Test{
		{
			cryptocurrency: "Bitcoin",
			vote:           1,
		},
		{
			cryptocurrency: "Bitcoin",
			vote:           1,
		},
		{
			cryptocurrency: "Bitcoin",
			vote:           -1,
		},
		{
			cryptocurrency: "Bitcoin",
			vote:           1,
		},
		{
			cryptocurrency: "Iota",
			vote:           1,
		},
		{
			cryptocurrency: "Iota",
			vote:           1,
		},
		{
			cryptocurrency: "Ethereum",
			vote:           1,
		},
		{
			cryptocurrency: "Ethereum",
			vote:           -1,
		},
		{
			cryptocurrency: "Ethereum",
			vote:           -1,
		},
	}

	for _, x := range s {
		if !checkTypes(x.cryptocurrency, x.vote) {
			t.Error("Wrong type to validate vote")
		}
	}
}

func checkTypes(a string, b int) bool {
	if reflect.TypeOf(a).String() == "string" && reflect.TypeOf(b).String() == "int" {
		return true
	} else {
		return false
	}
}
