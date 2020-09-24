package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultClient *mongoClient

type Config interface {
	GetMongodbConnectionString() string
	GetMongodbDatabase() string
}

type mongoClient struct {
	client *mongo.Client
	config Config
}

func InitDefaultClient(config Config) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetMongodbConnectionString()))
	if err != nil {
		return nil, err
	}
	defaultClient = &mongoClient{
		client: c,
		config: config,
	}
	return func() {
		if err = c.Disconnect(ctx); err != nil {
			panic(err)
		}
	}, nil
}

func Devcluster() *mongo.Database {
	return defaultClient.client.Database(defaultClient.config.GetMongodbDatabase())
}

func ClusterRequests() *mongo.Collection {
	return Devcluster().Collection("clusterRequests")
}

func Clusters() *mongo.Collection {
	return Devcluster().Collection("clusters")
}

func Users() *mongo.Collection {
	return Devcluster().Collection("users")
}
