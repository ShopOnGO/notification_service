package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"notification/configs"
)

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongo(conf *configs.Config) *Mongo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Mongo.URI))
	if err != nil {
		panic(err)
	}

	db := client.Database(conf.Mongo.Database)

	return &Mongo{
		Client:   client,
		Database: db,
	}
}
