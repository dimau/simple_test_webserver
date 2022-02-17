package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID `bson:"_id"`
	Isbn   string             `bson:"isbn"`
	Title  string             `bson:"title"`
	Author string             `bson:"author"`
	Price  float64            `bson:"price"`
}
