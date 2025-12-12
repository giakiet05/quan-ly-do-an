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
	GetByEmailAndClassroom(ctx context.Context, email string, classroomID string) (*model.ClassroomInvitation, error)
	GetPendingByEmail(ctx context.Context, email string) ([]model.ClassroomInvitation, error)
	GetPendingByUserID(ctx context.Context, userID string) ([]model.ClassroomInvitation, error)
	GetByClassroom(ctx context.Context, classroomID string, status *model.InvitationStatus, page, pageSize int) ([]model.ClassroomInvitation, int64, error)
	Update(ctx context.Context, invitation *model.ClassroomInvitation) (*model.ClassroomInvitation, error)
	UpdateStatus(ctx context.Context, invitationID string, status model.InvitationStatus) error
	Delete(ctx context.Context, invitationID string) error
}

type classroomInvitationRepo struct {
	collection *mongo.Collection
}

func NewClassroomInvitationRepo(db *mongo.Database) ClassroomInvitationRepo {
	return &classroomInvitationRepo{
		collection: db.Collection(config.ClassroomInvitationColName),
	}
}

func (r *classroomInvitationRepo) Create(ctx context.Context, invitation *model.ClassroomInvitation) (*model.ClassroomInvitation, error) {
	invitation.CreatedAt = time.Now()

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
	objectID, err := primitive.ObjectIDFromHex(invitationID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var invitation model.ClassroomInvitation
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invitation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *classroomInvitationRepo) GetByEmailAndClassroom(ctx context.Context, email string, classroomID string) (*model.ClassroomInvitation, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var invitation model.ClassroomInvitation
	filter := bson.M{
		"email":        email,
		"classroom_id": classroomObjectID,
		"status":       model.InvitationStatusPending,
	}

	err = r.collection.FindOne(ctx, filter).Decode(&invitation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *classroomInvitationRepo) GetPendingByEmail(ctx context.Context, email string) ([]model.ClassroomInvitation, error) {
	filter := bson.M{
		"email":  email,
		"status": model.InvitationStatusPending,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invitations []model.ClassroomInvitation
	if err := cursor.All(ctx, &invitations); err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *classroomInvitationRepo) GetPendingByUserID(ctx context.Context, userID string) ([]model.ClassroomInvitation, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := bson.M{
		"invited_to": objectID,
		"status":     model.InvitationStatusPending,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invitations []model.ClassroomInvitation
	if err := cursor.All(ctx, &invitations); err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *classroomInvitationRepo) GetByClassroom(ctx context.Context, classroomID string, status *model.InvitationStatus, page, pageSize int) ([]model.ClassroomInvitation, int64, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, 0, apperror.ErrInvalidID
	}

	filter := bson.M{"classroom_id": classroomObjectID}
	if status != nil {
		filter["status"] = *status
	}

	skip := int64((page - 1) * pageSize)

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
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

func (r *classroomInvitationRepo) Update(ctx context.Context, invitation *model.ClassroomInvitation) (*model.ClassroomInvitation, error) {
	invitation.UpdatedAt = time.Now()

	filter := bson.M{"_id": invitation.ID}
	update := bson.M{"$set": invitation}

	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updated model.ClassroomInvitation
	if err := result.Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrInvitationNotFound
		}
		return nil, err
	}

	return &updated, nil
}

func (r *classroomInvitationRepo) UpdateStatus(ctx context.Context, invitationID string, status model.InvitationStatus) error {
	objectID, err := primitive.ObjectIDFromHex(invitationID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrInvitationNotFound
	}

	return nil
}

func (r *classroomInvitationRepo) Delete(ctx context.Context, invitationID string) error {
	objectID, err := primitive.ObjectIDFromHex(invitationID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return apperror.ErrInvitationNotFound
	}

	return nil
}
