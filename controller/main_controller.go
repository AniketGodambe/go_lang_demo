package controller

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection constants
const connectionString = "mongodb+srv://aniketgodambe:aniketgodambe@cluster0.ehh0w.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbName = "contactdb"
const collName = "contacts"

// Global variable for the MongoDB collection
var Collection *mongo.Collection

// Initialize MongoDB connection
func InitDB() {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping to ensure connection is successful
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")
	Collection = client.Database(dbName).Collection(collName)
}
