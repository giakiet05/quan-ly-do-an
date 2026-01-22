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

	GetByChannelID(ctx context.Context, channelID string) (*model.Classroom, error)
	// Student management
	AddStudent(ctx context.Context, classroomID, studentID string) error
	RemoveStudent(ctx context.Context, classroomID, studentID string) error
	IsStudentInClassroom(ctx context.Context, classroomID, studentID string) (bool, error)
	GetStudentIDsByClassroomID(ctx context.Context, classroomID string) ([]string, error)

	// CoLecturer management
	AddCoLecturer(ctx context.Context, classroomID, coLecturerID string) error
	RemoveCoLecturer(ctx context.Context, classroomID, coLecturerID string) error
	IsCoLecturer(ctx context.Context, classroomID, userID string) (bool, error)
	IsLecturerOrCoLecturer(ctx context.Context, classroomID, userID string) (bool, error)

	// Classroom Join Code
	GetByInvitationCode(ctx context.Context, code string) (*model.Classroom, error)
	StudentCodeExistsInClassroom(ctx context.Context, classroomID, studentCode string) (bool, error)
	UpdateWhitelistStudentCode(ctx context.Context, classroomID string, entries []model.WhitelistEntry) error
	AddToWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error
	RemoveFromWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error
	ClearWhitelistStudentCode(ctx context.Context, classroomID string) error
	MarkWhitelistEntryAsJoined(ctx context.Context, classroomID, studentCode string, userID primitive.ObjectID) error
	ResetWhitelistEntry(ctx context.Context, classroomID, studentCode string) error
	IsStudentCodeClaimedInWhitelist(ctx context.Context, classroomID, studentCode string) (bool, error)
	RegenerateInvitationCode(ctx context.Context, classroomID, newCode string) error

	// Stats
	//GetProjectStats(ctx context.Context, classroomID, roundID, projectID string) (*model.ProjectStats, error)
	//GetRoundStats(ctx context.Context, classroomID, roundID string) (*model.RoundStats, error)

	IsLecturer(ctx context.Context, classroomID, lecturerID string) (bool, error)
	IsMember(ctx context.Context, classroomID, userID string) (bool, error)
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
func (c *classroomRepo) GetByChannelID(ctx context.Context, channelID string) (*model.Classroom, error) {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return nil, err
	}
	var classroom model.Classroom
	err = c.collection.FindOne(ctx, bson.M{"general_channel_id": channelObjectID}).Decode(&classroom)
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
		"lecturer_id": lecturerObjectID,
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
	if err := cursor.All(ctx, &classrooms); err != nil {
		return nil, 0, err
	}

	return classrooms, total, nil
}

func (c *classroomRepo) GetByStudent(ctx context.Context, studentID string, page, pageSize int) ([]model.Classroom, int64, error) {
	studentObjectID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{
		"student_ids": studentObjectID,
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
	if err := cursor.All(ctx, &classrooms); err != nil {
		return nil, 0, err
	}

	return classrooms, total, nil
}

func (c *classroomRepo) AddStudent(ctx context.Context, classroomID, studentID string) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	studentObjectID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$addToSet": bson.M{"student_ids": studentObjectID}}

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
	update := bson.M{"$pull": bson.M{"student_ids": studentObjectID}}

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
		"_id":         classroomObjectID,
		"student_ids": studentObjectID,
	}

	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *classroomRepo) GetStudentIDsByClassroomID(ctx context.Context, classroomID string) ([]string, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var result struct {
		StudentIDs []primitive.ObjectID `bson:"student_ids"`
	}

	err = c.collection.FindOne(
		ctx,
		bson.M{"_id": classroomObjectID},
		options.FindOne().SetProjection(bson.M{"student_ids": 1}),
	).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrClassroomNotFound
		}
		return nil, err
	}

	studentIDs := make([]string, 0, len(result.StudentIDs))
	for _, studentID := range result.StudentIDs {
		studentIDs = append(studentIDs, studentID.Hex())
	}

	return studentIDs, nil
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
		"_id":         classroomObjectID,
		"lecturer_id": lecturerObjectID,
	}

	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *classroomRepo) IsMember(ctx context.Context, classroomID, userID string) (bool, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id": classroomObjectID,
		"$or": []bson.M{
			{"lecturer._id": userObjectID},
			{"students._id": userObjectID},
		},
	}

	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ==================== CO-LECTURER MANAGEMENT ====================

func (c *classroomRepo) AddCoLecturer(ctx context.Context, classroomID, coLecturerID string) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	coLecturerObjectID, err := primitive.ObjectIDFromHex(coLecturerID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$addToSet": bson.M{"co_lecturer_ids": coLecturerObjectID}}

	result, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrClassroomNotFound
	}

	return nil
}

func (c *classroomRepo) RemoveCoLecturer(ctx context.Context, classroomID, coLecturerID string) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	coLecturerObjectID, err := primitive.ObjectIDFromHex(coLecturerID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$pull": bson.M{"co_lecturer_ids": coLecturerObjectID}}

	result, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrClassroomNotFound
	}

	return nil
}

func (c *classroomRepo) IsCoLecturer(ctx context.Context, classroomID, userID string) (bool, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id":             classroomObjectID,
		"co_lecturer_ids": userObjectID,
	}

	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *classroomRepo) IsLecturerOrCoLecturer(ctx context.Context, classroomID, userID string) (bool, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id": classroomObjectID,
		"$or": []bson.M{
			{"lecturer_id": userObjectID},
			{"co_lecturer_ids._id": userObjectID},
		},
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
	// This method is no longer needed since we don't store student_code in classroom
	// Student code validation should be done through user lookup instead
	// Keeping this as a placeholder or can be removed
	return false, nil
}

// UpdateWhitelistStudentCode replaces the entire whitelist with new entries
func (c *classroomRepo) UpdateWhitelistStudentCode(ctx context.Context, classroomID string, entries []model.WhitelistEntry) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$set": bson.M{"whitelist_student_code": entries}},
	)

	return err
}

// AddToWhitelistStudentCode adds student codes to the whitelist (no duplicates)
func (c *classroomRepo) AddToWhitelistStudentCode(ctx context.Context, classroomID string, studentCodes []string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Convert strings to WhitelistEntry structs
	entries := make([]model.WhitelistEntry, len(studentCodes))
	for i, code := range studentCodes {
		entries[i] = model.WhitelistEntry{StudentCode: code}
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{"_id": classroomOID},
		bson.M{"$addToSet": bson.M{"whitelist_student_code": bson.M{"$each": entries}}},
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
		bson.M{"$pull": bson.M{"whitelist_student_code": bson.M{"student_code": bson.M{"$in": studentCodes}}}},
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
		bson.M{"$set": bson.M{"whitelist_student_code": []model.WhitelistEntry{}}},
	)

	return err
}

// MarkWhitelistEntryAsJoined marks a whitelist entry as claimed by a user
func (c *classroomRepo) MarkWhitelistEntryAsJoined(ctx context.Context, classroomID, studentCode string, userID primitive.ObjectID) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	now := time.Now()
	_, err = c.collection.UpdateOne(ctx,
		bson.M{
			"_id":                                 classroomOID,
			"whitelist_student_code.student_code": studentCode,
		},
		bson.M{"$set": bson.M{
			"whitelist_student_code.$.joined_by": userID,
			"whitelist_student_code.$.joined_at": now,
		}},
	)

	return err
}

// ResetWhitelistEntry resets a whitelist entry by clearing joined_by and joined_at
func (c *classroomRepo) ResetWhitelistEntry(ctx context.Context, classroomID, studentCode string) error {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	_, err = c.collection.UpdateOne(ctx,
		bson.M{
			"_id":                                 classroomOID,
			"whitelist_student_code.student_code": studentCode,
		},
		bson.M{"$set": bson.M{
			"whitelist_student_code.$.joined_by": nil,
			"whitelist_student_code.$.joined_at": nil,
		}},
	)

	return err
}

// IsStudentCodeClaimedInWhitelist checks if a student code has already been claimed
func (c *classroomRepo) IsStudentCodeClaimedInWhitelist(ctx context.Context, classroomID, studentCode string) (bool, error) {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	count, err := c.collection.CountDocuments(ctx, bson.M{
		"_id": classroomOID,
		"whitelist_student_code": bson.M{
			"$elemMatch": bson.M{
				"student_code": studentCode,
				"joined_by":    bson.M{"$ne": nil},
			},
		},
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
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
