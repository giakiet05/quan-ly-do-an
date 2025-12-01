package repo

import (
	"context"
	"errors"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepo interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	GetByID(ctx context.Context, groupID string) (*model.Group, error)
	GetFilter(ctx context.Context, classroomID string, projectID *string, memberID *string) ([]model.Group, error)
	Update(ctx context.Context, filter Filter, update UpdateDocument) error
	Replace(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, groupID string) error

	IsLeader(ctx context.Context, groupID string, userID string) (bool, error)
	IsMember(ctx context.Context, groupID string, userID string) (bool, error)
}

type groupRepo struct {
	groupCollection *mongo.Collection
}

func NewGroupRepo(db *mongo.Database) GroupRepo {
	return &groupRepo{groupCollection: db.Collection(config.GroupColName)}
}

func (g *groupRepo) Create(ctx context.Context, group *model.Group) (*model.Group, error) {
	result, err := g.groupCollection.InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		group.ID = oid
	}

	return group, nil
}

func (g *groupRepo) GetByID(ctx context.Context, groupID string) (*model.Group, error) {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}

	var group model.Group
	err = g.groupCollection.FindOne(ctx, bson.M{"_id": groupOID}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (g *groupRepo) GetFilter(ctx context.Context, classroomID string, projectID *string, memberID *string) ([]model.Group, error) {
	filter := bson.M{}
	filter["classroom_id"] = classroomID

	if projectID != nil {
		filter["project_id"] = *projectID
	}

	if memberID != nil {
		filter["member_id"] = *memberID
	}

	cursor, err := g.groupCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var results []model.Group
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (g *groupRepo) Update(ctx context.Context, filter Filter, update UpdateDocument) error {
	_, err := g.groupCollection.UpdateMany(ctx, bson.M(filter), bson.M(update))
	return err
}

func (g *groupRepo) Replace(ctx context.Context, group *model.Group) error {
	res, err := g.groupCollection.ReplaceOne(ctx, bson.M{"_id": group.ID}, group)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return apperror.ErrNoFieldsToUpdate
	}

	return nil
}

func (g *groupRepo) Delete(ctx context.Context, groupID string) error {
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	res, err := g.groupCollection.DeleteOne(ctx, bson.M{"_id": groupObjectID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return apperror.ErrGroupNotFound
	}

	return nil
}

func (g *groupRepo) IsLeader(ctx context.Context, groupID string, userID string) (bool, error) {
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": groupObjectID, "leader_id": userObjectID}

	count, err := g.groupCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (g *groupRepo) IsMember(ctx context.Context, groupID string, userID string) (bool, error) {
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	// Find the group and filter only the members array
	var group model.Group
	err = g.groupCollection.FindOne(ctx, bson.M{"_id": groupObjectID, "members._id": userObjectID}).Decode(&group)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
