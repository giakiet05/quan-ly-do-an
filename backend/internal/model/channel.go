package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	AdminIDs     []primitive.ObjectID `bson:"admin_ids,omitempty" json:"admin_ids,omitempty"`
	Members      []UserInfo           `bson:"members" json:"members"`
	UserSettings []ChannelUserSetting `bson:"user_settings" json:"user_settings"`
	Background   *string              `bson:"background" json:"background"`
	Status       ChannelStatus        `bson:"status" json:"status"`
	Setting      ChannelSetting       `bson:"setting" json:"setting"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
}

type ChannelStatus string

const (
	ChannelStatusActive ChannelStatus = "active"
	ChannelStatusBlock  ChannelStatus = "block"
)

type ChannelUserSetting struct {
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Notification    bool               `bson:"notification" json:"notification"`
	TypingIndicator bool               `bson:"typing_indicator" json:"typing_indicator"`
	IsDeleted       bool               `bson:"is_deleted" json:"is_deleted"`
}

type ChannelSetting struct {
	AllowMemberMessage bool `bson:"allow_member_message" json:"allow_member_message"`
	AllowMedia         bool `bson:"allow_media" json:"allow_media"`
	AllowAttachments   bool `bson:"allow_attachments" json:"allow_attachments"`
	AllowPinMessages   bool `bson:"allow_pin_messages" json:"allow_pin_messages"`
	AllowMentions      bool `bson:"allow_mentions" json:"allow_mentions"`
}
