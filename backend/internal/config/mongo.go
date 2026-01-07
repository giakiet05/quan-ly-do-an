package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	db     *mongo.Database
)

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI") // e.g. mongodb://user:pass@localhost:27017
	if uri == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set")
	}

	db = client.Database(dbName)

	// Auto-create required collections if they don't exist
	if err := ensureCollections(ctx, db); err != nil {
		log.Fatalf("Collection initialization failed: %v", err)
	}

	log.Printf("Using database: %s\n", dbName)
	return client
}

func ensureCollections(ctx context.Context, db *mongo.Database) error {
	collections, err := db.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	required := []string{
		UserColName,
		ClassroomColName,
		ClassroomJoinRequestColName,
		ClassPostColName,
		GroupColName,
		ChannelColName,
		MessageColName,
		NotificationColName,
		EmailVerificationColName,
		PasswordResetColName,
	}

	existing := make(map[string]bool, len(collections))
	for _, c := range collections {
		existing[c] = true
	}

	for _, name := range required {
		if !existing[name] {
			// Create collection
			if err := db.CreateCollection(ctx, name); err != nil {
				return fmt.Errorf("failed to create collection %q: %w", name, err)
			}
			log.Printf("Created collection: %s", name)
		} else {
			log.Printf("Collection exists: %s", name)
		}
	}

	log.Println("All required collections are ready")
	return nil
}
