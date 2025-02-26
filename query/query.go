package query

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Operator string

type Query interface {
	Filter() bson.M
}

type emptyQuery struct{}

func (q *emptyQuery) Filter() bson.M {
	return bson.M{}
}

func Empty() Query {
	return &emptyQuery{}
}

func ID(id primitive.ObjectID) Query {
	return &Comparison{
		Operator: Equal,
		Field:    "_id",
		Value:    id,
	}
}