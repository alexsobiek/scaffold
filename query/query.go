package query

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Query interface {
	Filter() bson.M
}

func ID(id primitive.ObjectID) Query {
	return &Comparison{
		Operator: Equal,
		Field:    "_id",
		Value:    id,
	}
}
