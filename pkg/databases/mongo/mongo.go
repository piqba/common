package mgo

import (
	"context"
	"sync"

	"github.com/piqba/common/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

// NewMongoDbOnce return an unique instance of mongodb server
func NewMongoDbOnce(uri string) (*mongo.Client, error) {
	if uri == "" {
		uri = config.LoadEnvOrFallback(
			"APP_DATASOURCE_MONGODB_URL",
			"mongodb://localhost:27017",
		)
	}
	var clientInstance *mongo.Client
	var err error
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(uri)
		// Connect to MongoDB
		client, errConn := mongo.Connect(context.TODO(), clientOptions)
		if errConn != nil {
			err = errConn
		}
		// Check the connection
		errPing := client.Ping(context.TODO(), nil)
		if err != nil {

			err = errPing
		}
		clientInstance = client
	})
	return clientInstance, err
}

// NewMongoDb return an instance of mongodb server
func NewMongoDb(uri string) (*mongo.Client, error) {
	if uri == "" {
		uri = config.LoadEnvOrFallback(
			"APP_DATASOURCE_MONGODB_URL",
			"mongodb://localhost:27017",
		)
	}
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, err
}
