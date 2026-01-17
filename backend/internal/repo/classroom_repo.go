package repo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClassroomRepo interface {
	// Classroom CRUD
	Create(ctx context.Context, classroom *model.Classroom) (*model.Classroom, error)
	Update(ctx context.Context, classroom *model.Classroom) (*model.Classroom, error)
	Delete(ctx context.Context, classroomID string) error
	GetByID(ctx context.Context, classroomID string) (*model.Classroom, error)
	GetByUniversity(ctx context.Context, universityID string, page, pageSize int) ([]model.Classroom, int64, error)
	GetByLecturer(ctx context.Context, lecturerID string, page, pageSize int) ([]model.Classroom, int64, error)
	GetByStudent(ctx context.Context, studentID string, page, pageSize int) ([]model.Classroom, int64, error)

	// Student management
	AddStudent(ctx context.Context, classroomID string, student model.UserInfo) error
	RemoveStudent(ctx context.Context, classroomID, studentID string) error
	IsStudentInClassroom(ctx context.Context, classroomID, studentID string) (bool, error)

	// Classroom Join Code
	GetByInvitationCode(ctx context.Context, code string) (*model.Classroom, error)
	StudentCodeExistsInClassroom(ctx context.Context, classroomID, studentCode string) (bool, error)
	UpdateWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error
	AddToWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error
	RemoveFromWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error
	ClearWhitelistStudentCode(ctx context.Context, classroomID string) error
	RegenerateInvitationCode(ctx context.Context, classroomID, newCode string) error

	// Rounds
	CreateRound(ctx context.Context, classroomID string, round *model.ProjectRound) error
	GetRound(ctx context.Context, classroomID, roundID string) (*model.ProjectRound, error)
	ListRounds(ctx context.Context, classroomID string, page, pageSize int) ([]model.ProjectRound, int64, error)

	// Stats
	//GetProjectStats(ctx context.Context, classroomID, roundID, projectID string) (*model.ProjectStats, error)
	//GetRoundStats(ctx context.Context, classroomID, roundID string) (*model.RoundStats, error)

	IsLecturer(ctx context.Context, classroomID, lecturerID string) (bool, error)
}

type classroomRepo struct {
	collection *mongo.Collection
}

func NewClassroomRepo(db *mongo.Database) ClassroomRepo {
	return &classroomRepo{collection: db.Collection(config.ClassroomColName)}
}

// ==================== CLASSROOM ====================
func (c *classroomRepo) Create(ctx context.Context, classroom *model.Classroom) (*model.Classroom, error) {
	classroom.CreatedAt = time.Now()
	result, err := c.collection.InsertOne(ctx, classroom)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		classroom.ID = oid
	}
	return classroom, nil
}

func (c *classroomRepo) GetByID(ctx context.Context, classroomID string) (*model.Classroom, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}

	var classroom model.Classroom
	err = c.collection.FindOne(ctx, bson.M{"_id": classroomObjectID}).Decode(&classroom)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrClassroomNotFound
		}
		return nil, err
	}
	return &classroom, nil
}

// Update updates an existing classroom
func (c *classroomRepo) Update(ctx context.Context, classroom *model.Classroom) (*model.Classroom, error) {
	filter := bson.M{"_id": classroom.ID}
	update := bson.M{"$set": classroom}

	result, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, apperror.ErrClassroomNotFound
	}

	return classroom, nil
}

// Delete soft deletes a classroom (or hard delete if you prefer)
func (c *classroomRepo) Delete(ctx context.Context, classroomID string) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Hard delete
	result, err := c.collection.DeleteOne(ctx, bson.M{"_id": classroomObjectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return apperror.ErrClassroomNotFound
	}

	return nil
}

func (c *classroomRepo) GetByUniversity(ctx context.Context, universityID string, page, pageSize int) ([]model.Classroom, int64, error) {
	universityObjectID, err := primitive.ObjectIDFromHex(universityID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{"university_id": universityObjectID}
	skip := (page - 1) * pageSize

	// Count total
	total, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Paginated results
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := c.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var classrooms []model.Classroom
	if err := cursor.All(ctx, &classrooms); err != nil {
		return nil, 0, err
	}

	return classrooms, total, nil
}

func (c *classroomRepo) GetByLecturer(ctx context.Context, lecturerID string, page, pageSize int) ([]model.Classroom, int64, error) {
	lecturerObjectID, err := primitive.ObjectIDFromHex(lecturerID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{
		"lecturer._id": lecturerObjectID,
	}

	skip := int64((page - 1) * pageSize)

	total, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	cursor, err := c.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var classrooms []model.Classroom
	cursor.All(ctx, &classrooms)

	return classrooms, total, nil
}

func (c *classroomRepo) GetByStudent(ctx context.Context, studentID string, page, pageSize int) ([]model.Classroom, int64, error) {
	studentObjectID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{
		"students._id": studentObjectID,
	}

	log.Printf("🔍 GetByStudent - studentID: %s, filter: %+v", studentID, filter)

	skip := int64((page - 1) * pageSize)

	total, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("❌ CountDocuments error: %v", err)
		return nil, 0, err
	}

	log.Printf("📊 CountDocuments result: %d classrooms found", total)

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	cursor, err := c.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var classrooms []model.Classroom
	cursor.All(ctx, &classrooms)

	return classrooms, total, nil
}

// ==================== ROUNDS ====================
func (c *classroomRepo) CreateRound(ctx context.Context, classroomID string, round *model.ProjectRound) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}

	round.ID = primitive.NewObjectID()
	round.CreatedAt = time.Now()

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$push": bson.M{"rounds": round}}

	_, err = c.collection.UpdateOne(ctx, filter, update)
	return err
}

func (c *classroomRepo) GetRound(ctx context.Context, classroomID, roundID string) (*model.ProjectRound, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}
	roundObjectID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return nil, err
	}

	// Aggregation pipeline để lấy round cụ thể
	pipeline := []bson.M{
		{"$match": bson.M{"_id": classroomObjectID}},
		{"$unwind": "$rounds"},
		{"$match": bson.M{"rounds._id": roundObjectID}},
		{"$replaceRoot": bson.M{"newRoot": "$rounds"}},
	}

	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var round model.ProjectRound
	if !cursor.Next(ctx) {
		return nil, apperror.ErrRoundNotFound
	}
	cursor.Decode(&round)
	return &round, nil
}
func (c *classroomRepo) ListRounds(ctx context.Context, classroomID string, page, pageSize int) ([]model.ProjectRound, int64, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize

	pipeline := []bson.M{
		{"$match": bson.M{"_id": classroomObjectID}},
		{"$unwind": "$rounds"},
		{"$sort": bson.D{{Key: "rounds.created_at", Value: -1}}},
		{"$skip": skip},
		{"$limit": pageSize},
		{"$replaceRoot": bson.M{"newRoot": "$rounds"}},
	}

	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var rounds []model.ProjectRound
	if err := cursor.All(ctx, &rounds); err != nil {
		return nil, 0, err
	}

	// COUNT TOTAL (pipeline riêng KHÔNG có skip/limit)
	countPipeline := []bson.M{
		{"$match": bson.M{"_id": classroomObjectID}},
		{"$unwind": "$rounds"},
		{"$count": "total"},
	}

	countCursor, err := c.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0]["total"].(int64)
	}

	return rounds, total, nil
}

func (c *classroomRepo) AddStudent(ctx context.Context, classroomID string, student model.UserInfo) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$addToSet": bson.M{"students": student}}

	result, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrClassroomNotFound
	}

	return nil
}

func (c *classroomRepo) RemoveStudent(ctx context.Context, classroomID, studentID string) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	studentObjectID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$pull": bson.M{"students": bson.M{"_id": studentObjectID}}}

	result, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrClassroomNotFound
	}

	return nil
}

func (c *classroomRepo) IsStudentInClassroom(ctx context.Context, classroomID, studentID string) (bool, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	studentObjectID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id":          classroomObjectID,
		"students._id": studentObjectID,
	}

	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *classroomRepo) IsLecturer(ctx context.Context, classroomID, lecturerID string) (bool, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	lecturerObjectID, err := primitive.ObjectIDFromHex(lecturerID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id":          classroomObjectID,
		"lecturer._id": lecturerObjectID,
	}

	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ==================== CLASSROOM JOIN CODE ====================

func (c *classroomRepo) GetByInvitationCode(ctx context.Context, code string) (*model.Classroom, error) {
	var classroom model.Classroom
	err := c.collection.FindOne(ctx, bson.M{"invitation_code": code}).Decode(&classroom)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrCodeInvalid
		}
		return nil, err
	}
	return &classroom, nil
}

func (c *classroomRepo) StudentCodeExistsInClassroom(ctx context.Context, classroomID, studentCode string) (bool, error) {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	count, err := c.collection.CountDocuments(ctx, bson.M{
		"_id":                   classroomOID,
		"students.student_code": studentCode,
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateWhitelistStudentCode replaces the entire whitelist with new list
func (c *classroomRepo) UpdateWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$set": bson.M{"whitelist_student_code": studentCodes}},
	)

	return err
}

// AddToWhitelistStudentCode adds student codes to the whitelist (no duplicates)
func (c *classroomRepo) AddToWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$addToSet": bson.M{"whitelist_student_code": bson.M{"$each": studentCodes}}},
	)

	return err
}

// RemoveFromWhitelistStudentCode removes student codes from the whitelist
func (c *classroomRepo) RemoveFromWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$pull": bson.M{"whitelist_student_code": bson.M{"$in": studentCodes}}},
	)

	return err
}

// ClearWhitelistStudentCode clears all student codes from whitelist
func (c *classroomRepo) ClearWhitelistStudentCode(ctx context.Context, classroomID string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$set": bson.M{"whitelist_student_code": []string{}}},
	)

	return err
}

func (c *classroomRepo) RegenerateInvitationCode(ctx context.Context, classroomID, newCode string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$set": bson.M{"invitation_code": newCode}},
	)

	return err
}
