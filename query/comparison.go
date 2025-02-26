package query

import "go.mongodb.org/mongo-driver/bson"

type ComparisonOperator Operator

const (
	Equal              ComparisonOperator = "$eq"
	GreaterThan        ComparisonOperator = "$gt"
	GreaterThanOrEqual ComparisonOperator = "$gte"
	In                 ComparisonOperator = "$in"
	LessThan           ComparisonOperator = "$lt"
	LessThanOrEqual    ComparisonOperator = "$lte"
	NotEqual           ComparisonOperator = "$ne"
	NotIn              ComparisonOperator = "$nin"
)

type Comparison struct {
	Operator ComparisonOperator
	Field    string
	Value    any
}

func (q *Comparison) Filter() bson.M {
	return bson.M{q.Field: bson.M{string(q.Operator): q.Value}}
}
