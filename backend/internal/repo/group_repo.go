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
)

type GroupRepo interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	GetByID(ctx context.Context, groupID string) (*model.Group, error)
	GetFilter(ctx context.Context, classroomID string, projectID *string, memberID *string) ([]model.Group, error)
	Update(ctx context.Context, filter Filter, update UpdateDocument) error
	Replace(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, groupID string) error

	AddMember(ctx context.Context, groupID string, user *model.User) error
	RemoveMember(ctx context.Context, groupID string, requesterID string) error

	GetInvitationByID(ctx context.Context, groupID, invitationID primitive.ObjectID) (*model.JoinGroupInvitation, error)
	GetJoinRequestByID(ctx context.Context, groupID, requestID primitive.ObjectID) (*model.JoinGroupRequest, error)

	IsLeader(ctx context.Context, groupID string, userID string) (bool, error)
	IsMember(ctx context.Context, groupID string, userID string) (bool, error)
	IsMaxMemberReached(ctx context.Context, groupID string) (bool, error)
	ReportExistsByPeriod(ctx context.Context, classroomID string, periodID string) (bool, error)
	IsAlreadyInGroup(ctx context.Context, userID string) (bool, error)
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

func (g *groupRepo) GetFilter(
	ctx context.Context,
	classroomID string,
	projectID *string,
	memberID *string,
) ([]model.Group, error) {

	filter := bson.M{}

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}
	filter["classroom_id"] = classroomOID

	if projectID != nil {
		projectOID, err := primitive.ObjectIDFromHex(*projectID)
		if err != nil {
			return nil, err
		}
		filter["project_id"] = projectOID
	}

	if memberID != nil {
		memberOID, err := primitive.ObjectIDFromHex(*memberID)
		if err != nil {
			return nil, err
		}
		// adjust field name based on schema
		filter["members._id"] = memberOID
	}

	cursor, err := g.groupCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Group
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (g *groupRepo) Update(ctx context.Context, filter Filter, update UpdateDocument) error {
	result, err := g.groupCollection.UpdateMany(ctx, bson.M(filter), bson.M(update))
	if err != nil {
		return err
	}

	// Check if any documents were actually modified
	if result.ModifiedCount == 0 && result.MatchedCount == 0 {
		return apperror.ErrGroupNotFound
	}

	return nil
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

func (g *groupRepo) AddMember(ctx context.Context, groupID string, user *model.User) error {
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": groupObjectID,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"members": bson.M{
				"_id":   user.ID,
				"full_name": user.FullName,
				"avatar":    user.Avatar,
				"joined_at": time.Now(),
			},
		},
	}

	res, err := g.groupCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (g *groupRepo) RemoveMember(ctx context.Context, groupID string, requesterID string) error {
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": groupObjectID,
	}

	update := bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"user_id": requesterObjectID,
			},
		},
	}

	res, err := g.groupCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (g *groupRepo) GetInvitationByID(
	ctx context.Context,
	groupID, invitationID primitive.ObjectID,
) (*model.JoinGroupInvitation, error) {
	var result struct {
		JoinInvitation []model.JoinGroupInvitation `bson:"join_invitations"`
	}

	err := g.groupCollection.FindOne(
		ctx,
		bson.M{
			"_id":                  groupID,
			"join_invitations._id": invitationID,
		},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	for _, in := range result.JoinInvitation {
		if in.ID == invitationID {
			return &in, nil
		}
	}

	return nil, apperror.ErrNotFound
}

func (g *groupRepo) GetJoinRequestByID(
	ctx context.Context,
	groupID, requestID primitive.ObjectID,
) (*model.JoinGroupRequest, error) {
	var result struct {
		JoinRequest []model.JoinGroupRequest `bson:"join_requests"`
	}

	err := g.groupCollection.FindOne(
		ctx,
		bson.M{
			"_id":               groupID,
			"join_requests._id": requestID,
		},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	for _, re := range result.JoinRequest {
		if re.ID == requestID {
			return &re, nil
		}
	}

	return nil, apperror.ErrNotFound
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

func (g *groupRepo) IsMaxMemberReached(ctx context.Context, groupID string) (bool, error) {
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, err
	}

	var result struct {
		Members []struct{} `bson:"members"`
		Setting struct {
			MaxMembers int `bson:"max_members"`
		} `bson:"setting"`
	}

	err = g.groupCollection.FindOne(
		ctx,
		bson.M{"_id": groupObjectID},
	).Decode(&result)
	if err != nil {
		return false, err
	}

	if result.Setting.MaxMembers <= 0 {
		return false, nil
	}

	return len(result.Members) >= result.Setting.MaxMembers, nil
}

func (g *groupRepo) ReportExistsByPeriod(
	ctx context.Context,
	classroomID string,
	periodID string,
) (bool, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, err
	}

	periodOID, err := primitive.ObjectIDFromHex(periodID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"classroom_id":             classroomOID,
		"reports.report_period_id": periodOID,
	}

	count, err := g.groupCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (g *groupRepo) IsAlreadyInGroup(ctx context.Context, userID string) (bool, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"members._id": userOID,
	}

	err = g.groupCollection.FindOne(ctx, filter).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
