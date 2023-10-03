package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Products struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float32            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
}
