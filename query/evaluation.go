package query

import "go.mongodb.org/mongo-driver/bson"

type EvaluationOperator string

const (
	Expr  EvaluationOperator = "$expr"
	Json  EvaluationOperator = "$jsonSchema"
	Mod   EvaluationOperator = "$mod"
	Regex EvaluationOperator = "$regex"
	Text  EvaluationOperator = "$text"
	Where EvaluationOperator = "$where"
)

type Expression struct {
	Operator ComparisonOperator
	Value    any
}

func (q *Expression) Filter() bson.M {
	return bson.M{string(q.Operator): q.Value}
}

type JsonSchema struct {
	Schema bson.M
}

func (q *JsonSchema) Filter() bson.M {
	return bson.M{string(Json): q.Schema}
}

type Modulo struct {
	Field     string
	Divisor   int
	Remainder int
}

func (q *Modulo) Filter() bson.M {
	return bson.M{q.Field: bson.M{string(Mod): []any{q.Divisor, q.Remainder}}}
}

type RegexOption string

const (
	RegexCaseInsensitive RegexOption = "i"
	RegexMultiline       RegexOption = "m"
	RegexDotAll          RegexOption = "s"
	RegexExtended        RegexOption = "x"
	RegexUnicode         RegexOption = "u"
)

type RegularExpression struct {
	Field   string
	Pattern string
	Options []RegexOption
}

func (q *RegularExpression) Filter() bson.M {
	return bson.M{q.Field: bson.M{string(Regex): q.Pattern, "$options": q.Options}}
}

type TextSearch struct {
	Field              string
	Search             string
	Language           string
	CaseSensitive      bool
	DiacriticSensitive bool
}

func (q *TextSearch) Filter() bson.M {
	return bson.M{
		q.Field: bson.M{
			string(Text): bson.M{
				"$search":             q.Search,
				"$language":           q.Language,
				"$caseSensitive":      q.CaseSensitive,
				"$diacriticSensitive": q.DiacriticSensitive,
			},
		},
	}
}

type WhereSearch struct {
	Code string
}

func (q *WhereSearch) Filter() bson.M {
	return bson.M{string(Where): q.Code}
}
