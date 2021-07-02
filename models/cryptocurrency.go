package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cryptocurrency struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Initials    string             `json:"initials" bson:"initials"`
	Upvote      int32              `json:"upvote" bson:"upvote"`
	Downvote    int32              `json:"downvote" bson:"downvote"`
	Description string             `json:"description" bson:"description"`
}
