package repo

import (
	"context"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClassPostRepo interface {
	Create(ctx context.Context, post *model.ClassPost) (*model.ClassPost, error)
	GetByID(ctx context.Context, postID string) (*model.ClassPost, error)
	GetByClassroom(ctx context.Context, classroomID string, page, pageSize int) ([]model.ClassPost, int64, error)
	Update(ctx context.Context, post *model.ClassPost) (*model.ClassPost, error)
	Delete(ctx context.Context, postID string) error
	TogglePin(ctx context.Context, postID string, isPinned bool) error
}

type classPostRepo struct {
	collection *mongo.Collection
}

func NewClassPostRepo(db *mongo.Database) ClassPostRepo {
	return &classPostRepo{
		collection: db.Collection(config.ClassPostColName),
	}
}

func (r *classPostRepo) Create(ctx context.Context, post *model.ClassPost) (*model.ClassPost, error) {
	result, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	post.ID = result.InsertedID.(primitive.ObjectID)
	return post, nil
}

func (r *classPostRepo) GetByID(ctx context.Context, postID string) (*model.ClassPost, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	var post model.ClassPost
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *classPostRepo) GetByClassroom(ctx context.Context, classroomID string, page, pageSize int) ([]model.ClassPost, int64, error) {
	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, 0, apperror.ErrBadRequest
	}

	filter := bson.M{"classroom_id": classroomOID}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination, sort by pinned first, then created_at desc
	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSort(bson.D{{Key: "is_pinned", Value: -1}, {Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []model.ClassPost
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *classPostRepo) Update(ctx context.Context, post *model.ClassPost) (*model.ClassPost, error) {
	filter := bson.M{"_id": post.ID}
	update := bson.M{"$set": post}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, apperror.ErrPostNotFound
	}

	return post, nil
}

func (r *classPostRepo) Delete(ctx context.Context, postID string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return apperror.ErrPostNotFound
	}

	return nil
}

func (r *classPostRepo) TogglePin(ctx context.Context, postID string, isPinned bool) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	update := bson.M{"$set": bson.M{"is_pinned": isPinned}}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrPostNotFound
	}

	return nil
}
