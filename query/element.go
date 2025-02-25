package query

import "go.mongodb.org/mongo-driver/bson"

type ElementOperator string

const (
	Exists ElementOperator = "$exists"
	Type   ElementOperator = "$type"
)

type Element struct {
	Operator ElementOperator
	Value    any
}

func (q *Element) Filter() bson.M {
	return bson.M{string(q.Operator): q.Value}
}
