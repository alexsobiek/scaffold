package query

import "go.mongodb.org/mongo-driver/bson"

type ElementOperator Operator

const (
	Exists ElementOperator = "$exists"
	Type   ElementOperator = "$type"
)

type Element struct {
	Operator ElementOperator
	Field    string
	Value    any
}

func (q *Element) Filter() bson.M {
	return bson.M{q.Field: bson.M{string(q.Operator): q.Value}}
}
