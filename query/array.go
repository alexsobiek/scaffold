package query

import "go.mongodb.org/mongo-driver/bson"

type ArrayOperator Operator

const (
	All       ArrayOperator = "$all"
	ElemMatch ArrayOperator = "$elemMatch"
	Size      ArrayOperator = "$size"
)

type Array struct {
	Operator ArrayOperator
	Field    string
	Value    any
}

func (q *Array) Filter() bson.M {
	return bson.M{q.Field: bson.M{string(q.Operator): q.Value}}
}
