package query

import "go.mongodb.org/mongo-driver/bson"

type LogicalOperator Operator

const (
	And LogicalOperator = "$and"
	Not LogicalOperator = "$not"
	Nor LogicalOperator = "$nor"
	Or  LogicalOperator = "$or"
)

type Logical struct {
	Operator LogicalOperator
	Queries  []Query
}

func (q *Logical) Filter() bson.M {
	var filters []bson.M
	for _, query := range q.Queries {
		filters = append(filters, query.Filter())
	}
	return bson.M{string(q.Operator): filters}
}
