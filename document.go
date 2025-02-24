package scaffold

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document[T any] struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Created     primitive.DateTime `bson:"created" json:"created"`
	LastUpdated primitive.DateTime `bson:"last_updated" json:"last_updated"`
	Data        *T                 `bson:",inline" json:"document"`
	collection  *C[T]              `bson:"-"`
}

func createDocument[T any](collection *C[T], data T) *Document[T] {
	now := primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))
	return &Document[T]{
		ID:          primitive.NewObjectID(),
		Created:     now,
		LastUpdated: now,
		Data:        &data,
		collection:  collection,
	}
}

func (d *Document[T]) GetID() primitive.ObjectID {
	return d.ID
}

func (d *Document[T]) GetCreatedTime() primitive.DateTime {
	return d.Created
}

func (d *Document[T]) GetLastUpdateTime() primitive.DateTime {
	return d.LastUpdated
}

func (d *Document[T]) GetData() *T {
	return d.Data
}

// Set updates a single top-level field and triggers a DB update.
func (d *Document[T]) Set(ctx context.Context, field string, val interface{}) error {
	return d.SetMany(ctx, map[string]interface{}{field: val})
}

// SetMany updates multiple top-level fields and triggers a DB update.
func (d *Document[T]) SetMany(ctx context.Context, fields map[string]interface{}) error {
	// Create a map to track changed fields
	dbUpdates := bson.M{}

	// Create a lookup map for BSON tags to struct field names
	val := reflect.ValueOf(d.Data).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)

		names := getFieldNames(structField)

		var ok bool
		var value interface{}

		value, ok = fields[names.BsonField]

		if !ok {
			value, ok = fields[names.StructName]
		}

		if !ok {
			continue
		}

		field := val.FieldByName(names.StructName)

		// If still invalid, return an error
		if !field.IsValid() {
			return fmt.Errorf("field %s does not exist", names.StructName)
		}

		// Check if the type matches
		if field.Type() != reflect.TypeOf(value) {
			return fmt.Errorf("value type does not match field type for %s", names.StructName)
		}

		// Only update if the value is different from the current one

		if !reflect.DeepEqual(field.Interface(), value) {
			// Set the new value
			field.Set(reflect.ValueOf(value))

			// Add the changed field to the map with bson field name
			dbUpdates[names.BsonField] = value
		}
	}

	// Output the changed fields for demonstration
	if len(dbUpdates) > 0 {
		dbUpdates["last_updated"] = primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))

		_, err := d.collection.mc.UpdateByID(ctx, d.ID, bson.M{"$set": dbUpdates})
		return err
	}

	return nil
}

func (d *Document[T]) Delete(ctx context.Context) error {

	err := d.collection.delete(ctx, d.ID)

	if err != nil {
		return err
	}

	_, err = d.collection.mc.DeleteOne(ctx, bson.M{"_id": d.ID})
	return err
}
