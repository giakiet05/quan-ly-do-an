package main

import (
	"context"
	"log"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	log.Println("🚀 Starting classroom UserInfo to IDs migration...")

	// Load config from parent directory (.env in backend/)
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Find all classrooms that have old format (with lecturer, co_lecturers, students as UserInfo)
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "lecturer", Value: bson.D{{Key: "$exists", Value: true}}}},
					bson.D{{Key: "co_lecturers", Value: bson.D{{Key: "$exists", Value: true}}}},
					bson.D{{Key: "students", Value: bson.D{{Key: "$exists", Value: true}}}},
				}},
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

	log.Printf("📊 Found %d classrooms with old UserInfo format", len(classrooms))

	if len(classrooms) == 0 {
		log.Println("✅ No migration needed - all classrooms already have new format")
		return
	}

	// Migrate each classroom
	migratedCount := 0
	for _, classroom := range classrooms {
		classroomID := classroom["_id"]

		update := bson.M{}
		unset := bson.M{}

		// Convert lecturer UserInfo to lecturer_id
		if lecturer, ok := classroom["lecturer"].(bson.M); ok {
			if lecturerID, ok := lecturer["_id"].(primitive.ObjectID); ok {
				update["lecturer_id"] = lecturerID
				unset["lecturer"] = ""
			}
		}

		// Convert co_lecturers []UserInfo to co_lecturer_ids []ObjectID
		if coLecturers, ok := classroom["co_lecturers"].(bson.A); ok {
			coLecturerIDs := make([]primitive.ObjectID, 0, len(coLecturers))
			for _, coLecturer := range coLecturers {
				if coLecturerMap, ok := coLecturer.(bson.M); ok {
					if coLecturerID, ok := coLecturerMap["_id"].(primitive.ObjectID); ok {
						coLecturerIDs = append(coLecturerIDs, coLecturerID)
					}
				}
			}
			update["co_lecturer_ids"] = coLecturerIDs
			unset["co_lecturers"] = ""
		}

		// Convert students []UserInfo to student_ids []ObjectID
		if students, ok := classroom["students"].(bson.A); ok {
			studentIDs := make([]primitive.ObjectID, 0, len(students))
			for _, student := range students {
				if studentMap, ok := student.(bson.M); ok {
					if studentID, ok := studentMap["_id"].(primitive.ObjectID); ok {
						studentIDs = append(studentIDs, studentID)
					}
				}
			}
			update["student_ids"] = studentIDs
			unset["students"] = ""
		}

		// Skip if nothing to update
		if len(update) == 0 {
			log.Printf("⚠️  Skipping classroom %v - nothing to migrate", classroomID)
			continue
		}

		// Update the classroom
		filter := bson.M{"_id": classroomID}
		updateDoc := bson.M{}
		if len(update) > 0 {
			updateDoc["$set"] = update
		}
		if len(unset) > 0 {
			updateDoc["$unset"] = unset
		}

		result, err := classroomCollection.UpdateOne(ctx, filter, updateDoc)
		if err != nil {
			log.Printf("❌ Failed to update classroom %v: %v", classroomID, err)
			continue
		}

		if result.ModifiedCount > 0 {
			migratedCount++
			log.Printf("✅ Migrated classroom %v", classroomID)
		}
	}

	log.Printf("🎉 Migration completed! Migrated %d/%d classrooms", migratedCount, len(classrooms))
}
