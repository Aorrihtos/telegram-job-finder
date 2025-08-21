package db

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectToDb() (*mongo.Client, error) {
	// Create a new mongo client and connect to the server
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongopts := options.Client().ApplyURI(os.Getenv("MONGO_URL")).SetServerAPIOptions(serverAPI)
	cl, err := mongo.Connect(mongopts)
	if err != nil { return client, err }
	client = cl
	log.Println("Connected to MongoDB successfully")
	return cl, nil
}

func GetJobsCollection() *mongo.Collection {
	return client.Database("telegram-job-finder").Collection("jobs")
}

func GetUsersCollection() *mongo.Collection {
	return client.Database("telegram-job-finder").Collection("users")
}