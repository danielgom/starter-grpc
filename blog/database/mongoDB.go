package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	Collection *mongo.Collection
)

func Init() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Println("Connecting to MongoDB")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatalln("error connecting to the Mongo client", err.Error())
	}

	Collection = client.Database("blog_gRPC").Collection("blog")
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln("error pinging MongoDB instance", err.Error())
	}

	log.Println("Successfully connected to MongoDB")
}
