package main

import (
	"context"
	"log"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	log.Println("🚀 Starting ClassPost Author to ID migration...")

	// Load .env from multiple possible locations
	if err := godotenv.Load("../../.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("⚠️  .env file not found, using environment variables")
		}
	}
	config.LoadConfig()

	// Connect to MongoDB
	client := config.NewMongoClient()
	db := client.Database(config.Cfg.DBName)
	collection := db.Collection("class_posts")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Find all posts with old format (author as UserInfo)
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "author", Value: bson.D{{Key: "$exists", Value: true}}},
				{Key: "author._id", Value: bson.D{{Key: "$exists", Value: true}}},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("❌ Failed to find posts: %v", err)
	}
	defer cursor.Close(ctx)

	var posts []bson.M
	if err := cursor.All(ctx, &posts); err != nil {
		log.Fatalf("❌ Failed to decode posts: %v", err)
	}

	log.Printf("📊 Found %d posts with old author format", len(posts))

	if len(posts) == 0 {
		log.Println("✅ No migration needed - all posts already have new format")
		return
	}

	// Migrate each post
	migratedCount := 0
	for _, post := range posts {
		postID := post["_id"]
		author, ok := post["author"].(bson.M)
		if !ok {
			log.Printf("⚠️  Skipping post %v - invalid author format", postID)
			continue
		}

		authorID := author["_id"]

		// Update post
		filter := bson.M{"_id": postID}
		update := bson.M{
			"$set":   bson.M{"author_id": authorID},
			"$unset": bson.M{"author": ""},
		}

		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("❌ Failed to update post %v: %v", postID, err)
			continue
		}

		if result.ModifiedCount > 0 {
			migratedCount++
			log.Printf("✅ Migrated post %v", postID)
		}
	}

	log.Printf("🎉 Migration completed! Migrated %d/%d posts", migratedCount, len(posts))
}
