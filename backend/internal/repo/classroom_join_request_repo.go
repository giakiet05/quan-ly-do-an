package repo

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClassroomJoinRequestRepo interface {
	Create(ctx context.Context, req *model.ClassroomJoinRequest) error
	GetByID(ctx context.Context, id string) (*model.ClassroomJoinRequest, error)
	GetByUserAndClassroom(ctx context.Context, userID, classroomID string) (*model.ClassroomJoinRequest, error)
	GetPendingByClassroom(ctx context.Context, classroomID string, page, pageSize int) ([]model.ClassroomJoinRequest, int64, error)
	UpdateStatus(ctx context.Context, id string, status model.JoinRequestStatus, reviewerID *primitive.ObjectID) error
	CountPendingByClassroom(ctx context.Context, classroomID string) (int64, error)
}

type classroomJoinRequestRepo struct {
	collection *mongo.Collection
}

func NewClassroomJoinRequestRepo(db *mongo.Database) ClassroomJoinRequestRepo {
	return &classroomJoinRequestRepo{
		collection: db.Collection("classroom_join_requests"),
	}
}

func (r *classroomJoinRequestRepo) Create(ctx context.Context, req *model.ClassroomJoinRequest) error {
	_, err := r.collection.InsertOne(ctx, req)
	return err
}

func (r *classroomJoinRequestRepo) GetByID(ctx context.Context, id string) (*model.ClassroomJoinRequest, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrJoinRequestNotFound
	}

	var req model.ClassroomJoinRequest
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrJoinRequestNotFound
		}
		return nil, err
	}

	return &req, nil
}

func (r *classroomJoinRequestRepo) GetByUserAndClassroom(ctx context.Context, userID, classroomID string) (*model.ClassroomJoinRequest, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var req model.ClassroomJoinRequest
	err = r.collection.FindOne(ctx, bson.M{
		"user_id":      userOID,
		"classroom_id": classroomOID,
	}).Decode(&req)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not an error, just not found
		}
		return nil, err
	}

	return &req, nil
}

func (r *classroomJoinRequestRepo) GetPendingByClassroom(ctx context.Context, classroomID string, page, pageSize int) ([]model.ClassroomJoinRequest, int64, error) {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, 0, apperror.ErrInvalidID
	}

	filter := bson.M{
		"classroom_id": classroomOID,
		"status":       model.JoinRequestPending,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var requests []model.ClassroomJoinRequest
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (r *classroomJoinRequestRepo) UpdateStatus(ctx context.Context, id string, status model.JoinRequestStatus, reviewerID *primitive.ObjectID) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrInvalidID
	}

	now := time.Now()
	update := bson.M{
		"status":      status,
		"reviewed_at": now,
	}

	if reviewerID != nil {
		update["reviewed_by"] = reviewerID
	}

	_, err = r.collection.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": update},
	)

	return err
}

func (r *classroomJoinRequestRepo) CountPendingByClassroom(ctx context.Context, classroomID string) (int64, error) {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return 0, apperror.ErrInvalidID
	}

	return r.collection.CountDocuments(ctx, bson.M{
		"classroom_id": classroomOID,
		"status":       model.JoinRequestPending,
	})
}
