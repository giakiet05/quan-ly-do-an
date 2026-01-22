package repo

import (
	"context"
	"log"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChannelRepo interface {
	Create(ctx context.Context, req *dto.CreateChannelRequest, requesterID string) (*model.Channel, error)
	CreateClassroomChannel(ctx context.Context, adminIDs []string, students []model.UserInfo) (*model.Channel, error)
	CreateGroupChannel(ctx context.Context, leaderID string, members []model.UserInfo) (*model.Channel, error)
	GetByID(ctx context.Context, channelID string) (*model.Channel, error)
	GetByUserID(ctx context.Context, userID string, page int, pageSize int) ([]model.Channel, int64, error)
	GetByBothUserID(ctx context.Context, user1ID string, user2ID string) (*model.Channel, error)
	Update(ctx context.Context, channel *model.Channel) (*model.Channel, error)
	UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error
	Delete(ctx context.Context, channelID string, userID string) error
	IsMember(ctx context.Context, channelID string, userID string) (bool, error)
	AddMember(ctx context.Context, channelID string, userID string) error
	RemoveMember(ctx context.Context, channelID string, userID string) error
}

type channelRepo struct {
	channelCollection *mongo.Collection
}

func NewChannelRepo(db *mongo.Database) ChannelRepo {
	return &channelRepo{channelCollection: db.Collection(config.ChannelColName)}
}

func (c *channelRepo) Create(ctx context.Context, req *dto.CreateChannelRequest, requesterID string) (*model.Channel, error) {
	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	settings := make([]model.ChannelUserSetting, 0, len(req.Members))
	for _, m := range req.Members {
		settings = append(settings, model.ChannelUserSetting{
			UserID:          m.ID,
			Notification:    true,
			TypingIndicator: true,
			IsDeleted:       false,
		})
	}

	channel := &model.Channel{
		AdminIDs:     []primitive.ObjectID{requesterObjectID},
		Members:      req.Members,
		UserSettings: settings,
		Status:       model.ChannelStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result, err := c.channelCollection.InsertOne(ctx, channel)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		channel.ID = oid
	}

	return channel, nil
}

func (c *channelRepo) CreateClassroomChannel(
	ctx context.Context,
	adminIDs []string,
	students []model.UserInfo,
) (*model.Channel, error) {
	now := time.Now()

	userIDSet := make(map[primitive.ObjectID]bool)
	settings := make([]model.ChannelUserSetting, 0, len(students)+len(adminIDs))
	for _, m := range students {
		if !userIDSet[m.ID] {
			userIDSet[m.ID] = true
			settings = append(settings, model.ChannelUserSetting{
				UserID:          m.ID,
				Notification:    true,
				TypingIndicator: true,
				IsDeleted:       false,
			})
		}
	}

	adminObjectIDs := make([]primitive.ObjectID, 0, len(adminIDs))
	for _, id := range adminIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		adminObjectIDs = append(adminObjectIDs, oid)

		if !userIDSet[oid] {
			userIDSet[oid] = true
			settings = append(settings, model.ChannelUserSetting{
				UserID:          oid,
				Notification:    true,
				TypingIndicator: true,
				IsDeleted:       false,
			})
		}
	}

	channel := &model.Channel{
		ID:           primitive.NewObjectID(),
		AdminIDs:     adminObjectIDs,
		Members:      students,
		UserSettings: settings,
		Background:   nil,
		Status:       model.ChannelStatusActive,

		Setting: model.ChannelSetting{
			AllowMemberMessage: true,
			AllowMedia:         true,
			AllowAttachments:   true,
			AllowPinMessages:   false,
			AllowMentions:      true,
		},

		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := c.channelCollection.InsertOne(ctx, channel)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		channel.ID = oid
	}

	return channel, nil
}

func (c *channelRepo) CreateGroupChannel(ctx context.Context, leaderID string, members []model.UserInfo) (*model.Channel, error) {
	now := time.Now()

	settings := make([]model.ChannelUserSetting, 0, len(members))
	for _, m := range members {
		settings = append(settings, model.ChannelUserSetting{
			UserID:          m.ID,
			Notification:    true,
			TypingIndicator: true,
			IsDeleted:       false,
		})
	}

	leaderObjectID, err := primitive.ObjectIDFromHex(leaderID)
	if err != nil {
		return nil, err
	}
	adminIDs := []primitive.ObjectID{leaderObjectID}

	channel := &model.Channel{
		ID:           primitive.NewObjectID(),
		AdminIDs:     adminIDs,
		Members:      members,
		UserSettings: settings,
		Background:   nil,
		Status:       model.ChannelStatusActive,

		Setting: model.ChannelSetting{
			AllowMemberMessage: true,
			AllowMedia:         true,
			AllowAttachments:   true,
			AllowPinMessages:   false,
			AllowMentions:      true,
		},

		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := c.channelCollection.InsertOne(ctx, channel)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		channel.ID = oid
	}

	return channel, nil
}

func (c *channelRepo) GetByID(ctx context.Context, channelID string) (*model.Channel, error) {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return nil, err
	}

	var channel model.Channel
	err = c.channelCollection.FindOne(ctx, bson.M{"_id": channelObjectID}).Decode(&channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *channelRepo) GetByUserID(ctx context.Context, userID string, page int, pageSize int) ([]model.Channel, int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{
		"$or": []bson.M{
			{"members.user_id": userObjectID},
			{"admin_ids": userObjectID},
		},
	}
	opt := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.M{"updated_at": -1})

	cursor, err := c.channelCollection.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var channels []model.Channel
	if err = cursor.All(ctx, &channels); err != nil {
		return nil, 0, err
	}

	count, err := c.channelCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return channels, count, nil
}

func (c *channelRepo) GetByBothUserID(ctx context.Context, user1ID string, user2ID string) (*model.Channel, error) {
	log.Printf("🔍 GetByBothUserID called: user1=%s, user2=%s", user1ID, user2ID)

	user1ObjectID, err := primitive.ObjectIDFromHex(user1ID)
	if err != nil {
		log.Printf("❌ GetByBothUserID: Invalid user1ID: %v", err)
		return nil, err
	}

	user2ObjectID, err := primitive.ObjectIDFromHex(user2ID)
	if err != nil {
		log.Printf("❌ GetByBothUserID: Invalid user2ID: %v", err)
		return nil, err
	}

	// Find DM channel (exactly 2 members, both users are members)
	filter := bson.M{
		"$and": []bson.M{
			{"members.user_id": user1ObjectID},
			{"members.user_id": user2ObjectID},
		},
		"members": bson.M{"$size": 2}, // Ensure it's a DM channel (exactly 2 members)
	}

	log.Printf("🔎 GetByBothUserID filter: %+v", filter)

	var channel model.Channel
	err = c.channelCollection.FindOne(ctx, filter).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("📭 GetByBothUserID: No channel found (ErrChannelNotFound)")
			return nil, apperror.ErrChannelNotFound
		}
		log.Printf("❌ GetByBothUserID: Database error: %v", err)
		return nil, err
	}

	log.Printf("✅ GetByBothUserID: Found channel %s", channel.ID.Hex())
	return &channel, nil
}

func (c *channelRepo) Update(ctx context.Context, channel *model.Channel) (*model.Channel, error) {
	update := bson.M{
		"$set": bson.M{
			"settings":   channel.UserSettings,
			"status":     channel.Status,
			"updated_at": time.Now(),
		},
	}

	_, err := c.channelCollection.UpdateByID(ctx, channel.ID, update)
	if err != nil {
		return nil, err
	}

	var updated model.Channel
	err = c.channelCollection.FindOne(ctx, bson.M{"_id": channel.ID}).Decode(&updated)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (c *channelRepo) UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"members.user_id": userObjectID,
	}

	update := bson.M{
		"$set": bson.M{
			"members.$[elem].avatar": newAvatar,
			"updated_at":             time.Now(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": userObjectID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters).SetUpsert(false)

	_, err = c.channelCollection.UpdateMany(ctx, filter, update, opts)
	return err
}

func (c *channelRepo) Delete(ctx context.Context, channelID string, userID string) error {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"settings.$[elem].is_deleted": true,
			"updated_at":                  time.Now(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.user_id": userObjectID}},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	_, err = c.channelCollection.UpdateOne(ctx, bson.M{"_id": channelObjectID}, update, opts)
	return err
}

func (c *channelRepo) AddMember(ctx context.Context, channelID string, userID string) error {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$addToSet": bson.M{
			"members": bson.M{
				"user_id":   userObjectID,
				"joined_at": time.Now(),
			},
			"settings": bson.M{
				"user_id":          userObjectID,
				"notification":     true,
				"typing_indicator": true,
				"is_deleted":       false,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = c.channelCollection.UpdateOne(ctx, bson.M{"_id": channelObjectID}, update)
	return err
}

func (c *channelRepo) RemoveMember(ctx context.Context, channelID string, userID string) error {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"user_id": userObjectID,
			},
			"settings": bson.M{
				"user_id": userObjectID,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = c.channelCollection.UpdateOne(ctx, bson.M{"_id": channelObjectID}, update)
	return err
}

func (c *channelRepo) IsMember(ctx context.Context, channelID string, userID string) (bool, error) {
	channel, err := c.GetByID(ctx, channelID)
	if err != nil {
		return false, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	// Check if user is admin
	for _, adminID := range channel.AdminIDs {
		if adminID == userObjectID {
			return true, nil
		}
	}

	// Check if user is member
	for _, m := range channel.Members {
		if m.ID == userObjectID {
			return true, nil
		}
	}

	return false, nil
}
