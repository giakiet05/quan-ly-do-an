package repo

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClassroomInvitationRepo interface {
	Create(ctx context.Context, invitation *model.ClassroomInvitation) (*model.ClassroomInvitation, error)
	GetByID(ctx context.Context, invitationID string) (*model.ClassroomInvitation, error)
	GetByInviteeAndClassroom(ctx context.Context, inviteeID, classroomID string) (*model.ClassroomInvitation, error)
	GetPendingByInvitee(ctx context.Context, inviteeID string, page, pageSize int) ([]model.ClassroomInvitation, int64, error)
	GetPendingByClassroom(ctx context.Context, classroomID string, page, pageSize int) ([]model.ClassroomInvitation, int64, error)
	UpdateStatus(ctx context.Context, invitationID string, status model.InvitationStatus) error
	Delete(ctx context.Context, invitationID string) error
}

type classroomInvitationRepo struct {
	collection *mongo.Collection
}

func NewClassroomInvitationRepo(db *mongo.Database) ClassroomInvitationRepo {
	return &classroomInvitationRepo{collection: db.Collection(config.ClassroomInvitationColName)}
}

func (r *classroomInvitationRepo) Create(ctx context.Context, invitation *model.ClassroomInvitation) (*model.ClassroomInvitation, error) {
	invitation.CreatedAt = time.Now()
	invitation.Status = model.InvitationPending

	result, err := r.collection.InsertOne(ctx, invitation)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		invitation.ID = oid
	}

	return invitation, nil
}

func (r *classroomInvitationRepo) GetByID(ctx context.Context, invitationID string) (*model.ClassroomInvitation, error) {
	oid, err := primitive.ObjectIDFromHex(invitationID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var invitation model.ClassroomInvitation
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&invitation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *classroomInvitationRepo) GetByInviteeAndClassroom(ctx context.Context, inviteeID, classroomID string) (*model.ClassroomInvitation, error) {
	inviteeOID, err := primitive.ObjectIDFromHex(inviteeID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := bson.M{
		"invitee_id":   inviteeOID,
		"classroom_id": classroomOID,
		"status":       model.InvitationPending,
	}

	var invitation model.ClassroomInvitation
	err = r.collection.FindOne(ctx, filter).Decode(&invitation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *classroomInvitationRepo) GetPendingByInvitee(ctx context.Context, inviteeID string, page, pageSize int) ([]model.ClassroomInvitation, int64, error) {
	inviteeOID, err := primitive.ObjectIDFromHex(inviteeID)
	if err != nil {
		return nil, 0, apperror.ErrInvalidID
	}

	filter := bson.M{
		"invitee_id": inviteeOID,
		"status":     model.InvitationPending,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var invitations []model.ClassroomInvitation
	if err := cursor.All(ctx, &invitations); err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

func (r *classroomInvitationRepo) GetPendingByClassroom(ctx context.Context, classroomID string, page, pageSize int) ([]model.ClassroomInvitation, int64, error) {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, 0, apperror.ErrInvalidID
	}

	filter := bson.M{
		"classroom_id": classroomOID,
		"status":       model.InvitationPending,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var invitations []model.ClassroomInvitation
	if err := cursor.All(ctx, &invitations); err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

func (r *classroomInvitationRepo) UpdateStatus(ctx context.Context, invitationID string, status model.InvitationStatus) error {
	oid, err := primitive.ObjectIDFromHex(invitationID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"responded_at": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrInvitationNotFound
	}

	return nil
}

func (r *classroomInvitationRepo) Delete(ctx context.Context, invitationID string) error {
	oid, err := primitive.ObjectIDFromHex(invitationID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return apperror.ErrInvitationNotFound
	}

	return nil
}
