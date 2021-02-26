package database

import (
	"context"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)
var log = logging.MustGetLogger("main")

// CONSTANTS


func ListDatabases() []string {

	// SET UP A CLIENT
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://ccdc-scoreserver-database:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error occurred while connecting: %s", err)
	}
	defer client.Disconnect(ctx)

	// PING
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Error occurred while pinging: %s", err)
	}

	// LIST DATABASES
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("Error occurred while listing databases: %s", err)
	}
	return databases
}

func GetClient() (*mongo.Client, error ) {
	// SET UP A CLIENT
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://ccdc-scoreserver-database:27017/"))
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		//log.Fatalf("Error occurred while connecting: %s", err)
		return nil, err
	}
	defer client.Disconnect(ctx)

	// PING - MAKE SURE THE CONNECTION IS GOOD
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		//log.Fatalf("Error occurred while pinging: %s", err)
		return nil, err
	}

	return client, nil
}

func GetCollection(collection Collection) (*mongo.Collection, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.Database(databaseName).Collection(string(collection)), nil
}