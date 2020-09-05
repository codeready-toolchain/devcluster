package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDefaultClient(connectionString string) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	Client = c

	return func() {
		if err = c.Disconnect(ctx); err != nil {
			panic(err)
		}
	}, nil
}

func Devcluster() *mongo.Database {
	return Client.Database("devcluster")
}

func ClusterRequests() *mongo.Collection {
	return Devcluster().Collection("clusterRequests")
}

func Clusters() *mongo.Collection {
	return Devcluster().Collection("clusters")
}
