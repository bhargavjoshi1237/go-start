package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	log.Print("starting server...")

	// Initialize MongoDB connection
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://bhargavjoshi1237:Shiro123@cluster0.pwwbqv5.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatal(err)
	}

	// Set context with a timeout for MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Set up handlers
	http.HandleFunc("/", handler)
	http.HandleFunc("/xd", xdHandler)
	http.HandleFunc("/fetch", fetchHandler)

	// Start server
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!\n")
}

func xdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "xd")
}

// fetchHandler connects to the MongoDB collection, retrieves all data, and returns it as JSON.
func fetchHandler(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("jack").Collection("jack1")

	// Set context with a timeout for the MongoDB query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query all documents from the collection
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, "Failed to fetch data from MongoDB", http.StatusInternalServerError)
		log.Printf("Failed to fetch data: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Read all documents into a slice
	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		log.Printf("Failed to decode data: %v", err)
		return
	}

	// Convert the results to JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
