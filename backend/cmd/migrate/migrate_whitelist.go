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
	log.Println("🚀 Starting whitelist migration...")

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
	classroomCollection := db.Collection(config.ClassroomColName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find all classrooms that have old format whitelist (array of strings)
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "whitelist_student_code", Value: bson.D{{Key: "$exists", Value: true}}},
				{Key: "whitelist_student_code.0", Value: bson.D{{Key: "$type", Value: "string"}}},
			}},
		},
	}

	cursor, err := classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("❌ Failed to find classrooms: %v", err)
	}
	defer cursor.Close(ctx)

	var classrooms []bson.M
	if err := cursor.All(ctx, &classrooms); err != nil {
		log.Fatalf("❌ Failed to decode classrooms: %v", err)
	}

	log.Printf("📊 Found %d classrooms with old whitelist format", len(classrooms))

	if len(classrooms) == 0 {
		log.Println("✅ No migration needed - all classrooms already have new format")
		return
	}

	// Migrate each classroom
	migratedCount := 0
	for _, classroom := range classrooms {
		classroomID := classroom["_id"]
		oldWhitelist, ok := classroom["whitelist_student_code"].(bson.A)
		if !ok {
			log.Printf("⚠️  Skipping classroom %v - invalid whitelist format", classroomID)
			continue
		}

		// Convert old format to new format
		newWhitelist := make([]bson.M, 0, len(oldWhitelist))
		for _, code := range oldWhitelist {
			if codeStr, ok := code.(string); ok {
				newWhitelist = append(newWhitelist, bson.M{
					"student_code": codeStr,
					"joined_by":    nil,
					"joined_at":    nil,
				})
			}
		}

		// Update the classroom
		filter := bson.M{"_id": classroomID}
		update := bson.M{
			"$set": bson.M{
				"whitelist_student_code": newWhitelist,
			},
		}

		result, err := classroomCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("❌ Failed to update classroom %v: %v", classroomID, err)
			continue
		}

		if result.ModifiedCount > 0 {
			migratedCount++
			log.Printf("✅ Migrated classroom %v (%d codes)", classroomID, len(newWhitelist))
		}
	}

	log.Printf("🎉 Migration completed! Migrated %d/%d classrooms", migratedCount, len(classrooms))
}
