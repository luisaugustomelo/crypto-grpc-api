package crud

import (
	"context"

	mongo "klever/grpc/databases"
	"klever/grpc/models"
	system "klever/grpc/upvote/system"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateCryptocurrencyService(request *system.CreateCryptocurrencyRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.CreateCryptocurrencyResponse, error) {
	crypto := request.GetCrypto()

	name := crypto.GetName()
	description := crypto.GetDescription()
	initials := crypto.GetInitials()

	if name == "" || description == "" || initials == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty fields")
	}

	data := models.Cryptocurrency{
		Id:          primitive.NewObjectID(),
		Name:        name,
		Initials:    initials,
		Upvote:      0,
		Downvote:    0,
		Description: description,
	}

	findResult := db.FindOne(mongoCtx, bson.M{"name": name})

	isRegisteredCrypto := models.Cryptocurrency{}

	if err := findResult.Decode(&isRegisteredCrypto); err == nil {
		return nil, status.Error(codes.AlreadyExists, "Cryptocurrency already exists")
	}

	result, err := db.InsertOne(mongoCtx, data)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	crypto.Id = result.(primitive.ObjectID).Hex()

	response := &system.CreateCryptocurrencyResponse{Crypto: crypto}

	return response, nil
}

func UpdateCryptocurrencyService(request *system.UpdateCryptocurrencyRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.UpdateCryptocurrencyResponse, error) {

	crypto := request.GetCrypto()

	//Check if id is valid
	cryptoId, err := primitive.ObjectIDFromHex(crypto.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	name := crypto.GetName()
	description := crypto.GetDescription()
	initials := crypto.GetInitials()

	if name == "" || description == "" || initials == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty fields")
	}

	data := models.Cryptocurrency{
		Id:          cryptoId,
		Name:        name,
		Initials:    initials,
		Description: description,
	}

	result := db.FindOneAndUpdate(mongoCtx, bson.M{"_id": cryptoId}, bson.M{"$set": data})

	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.NotFound, "Cannot be find a crypto with this Object Id")
	}

	isUpdatedCrypto := models.Cryptocurrency{}

	err = result.Decode(&isUpdatedCrypto)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
	}

	response := &system.UpdateCryptocurrencyResponse{Crypto: crypto}

	return response, nil
}

func DeleteCryptocurrencyService(request *system.DeleteCryptocurrencyRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.DeleteCryptocurrencyResponse, error) {

	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := db.DeleteOne(mongoCtx, bson.M{"_id": cryptoId})

	if err != nil && result == 0 {
		return nil, status.Errorf(codes.NotFound, "Cannot be delete a crypto with this Object Id")
	}

	response := &system.DeleteCryptocurrencyResponse{Status: true}

	return response, nil
}

func ReadCryptocurrencyService(request *system.ReadCryptocurrencyRequest, db mongo.CollectionHelper, mongoCtx context.Context) (*system.ReadCryptocurrencyResponse, error) {
	cryptoId, err := primitive.ObjectIDFromHex(request.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result := db.FindOne(mongoCtx, bson.M{"_id": cryptoId})

	crypto := models.Cryptocurrency{}

	if err := result.Decode(&crypto); err == nil {
		status.Error(codes.NotFound, "Cannot be find a crypto with this Object Id")
	}

	response := &system.ReadCryptocurrencyResponse{
		Crypto: &system.Cryptocurrency{
			Id:          crypto.Id.Hex(),
			Name:        crypto.Name,
			Initials:    crypto.Initials,
			Upvote:      crypto.Upvote,
			Downvote:    crypto.Downvote,
			Description: crypto.Description,
		}}

	return response, nil
}
