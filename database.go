package scaffold

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
}

func NewDatabase(ctx context.Context, mongoUri string, database string) (*Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(database)

	return &Database{
		ctx:    ctx,
		client: client,
		db:     db,
	}, nil
}

func (d *Database) Collection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

func (d *Database) Close() {
	if err := d.client.Disconnect(d.ctx); err != nil {
		panic(err)
	}
}
