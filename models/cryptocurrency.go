package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cryptocurrency struct {
	id          primitive.ObjectID `json:"id" bson:"_id"`
	initials    string             `json:"initials" bson:"initials"`
	name        string             `json:"name" bson:"name"`
	description string             `json:"description" bson:"description"`
	upvote      int32              `json:"upvote" bson:"upvote"`
	downvote    int32              `json:"downvote" bson:"downvote"`
}
