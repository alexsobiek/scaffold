package scaffold

import (
	"reflect"
	"strings"
)

type bsonField struct {
	StructName string
	BsonField  string
}

func getFieldNames(f reflect.StructField) bsonField {
	bsonTag := f.Tag.Get("bson")

	if bsonTag != "" && bsonTag != "-" {
		bsonTag = strings.Split(bsonTag, ",")[0] // Extract only the actual BSON field name
	} else {
		bsonTag = f.Name
	}

	return bsonField{
		StructName: f.Name,
		BsonField:  bsonTag,
	}
}
