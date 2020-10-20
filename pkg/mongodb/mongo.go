package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultClient *mongoClient

type Config interface {
	GetMongodbConnectionString() string
	GetMongodbDatabase() string
	GetMongodbCA() string
}

type mongoClient struct {
	client *mongo.Client
	config Config
}

func InitDefaultClient(config Config) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := []*options.ClientOptions{
		options.Client().ApplyURI(config.GetMongodbConnectionString()),
	}
	if ca := config.GetMongodbCA(); ca != "" {
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM([]byte(ca)); !ok {
			return nil, errors.New("failed to parse the mongodb CA")
		}
		opts = append(opts, options.Client().SetTLSConfig(&tls.Config{RootCAs: roots}))
	}
	c, err := mongo.Connect(ctx, opts...)
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
