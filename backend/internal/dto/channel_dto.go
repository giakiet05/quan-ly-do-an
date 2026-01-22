package dto

import (
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

type CreateChannelRequest struct {
	Members []model.UserInfo `bson:"members" json:"members"`
}

type GetChannelByUserIDQuery struct {
	UserID   string `form:"user_id"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UpdateChannelRequest struct {
	ChannelID       string               `json:"channel_id"`
	Nickname        *string              `json:"nickname,omitempty"`
	Background      *string              `json:"background,omitempty"`
	Notification    *bool                `json:"notification,omitempty"`
	TypingIndicator *bool                `json:"typing_indicator,omitempty"`
	Status          *model.ChannelStatus `json:"status,omitempty"`
}

type ChannelResponse struct {
	ID                 string                   `json:"id"`
	AdminIDs           []string                 `json:"admin_ids,omitempty"`
	Members            []UserInfoResponse       `json:"members"`
	Settings           []ChannelSettingResponse `json:"settings"`
	Background         *string                  `json:"background"`
	Status             model.ChannelStatus      `json:"status"`
	UnreadMessageCount *int64                   `json:"unread_message_count,omitempty"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
}

type ChannelSettingResponse struct {
	UserID          string  `json:"user_id"`
	Nickname        *string `json:"nickname,omitempty"`
	Notification    bool    `json:"notification"`
	TypingIndicator bool    `json:"typing_indicator"`
}

func FromChannel(channel *model.Channel, unreadCount *int64) ChannelResponse {
	if channel == nil {
		return ChannelResponse{}
	}

	members := make([]UserInfoResponse, len(channel.Members))
	for i, m := range channel.Members {
		members[i] = UserInfoResponse{
			UserID:      m.ID.Hex(),
			FullName:    m.FullName,
			Email:       m.Email,
			Avatar:      m.Avatar,
			StudentCode: m.StudentCode,
		}
	}

	adminIDs := make([]string, len(channel.AdminIDs))
	for i, id := range channel.AdminIDs {
		adminIDs[i] = id.Hex()
	}

	settings := make([]ChannelSettingResponse, len(channel.UserSettings))
	for i, s := range channel.UserSettings {
		settings[i] = ChannelSettingResponse{
			UserID:          s.UserID.Hex(),
			Notification:    s.Notification,
			TypingIndicator: s.TypingIndicator,
		}
	}

	return ChannelResponse{
		ID:                 channel.ID.Hex(),
		AdminIDs:           adminIDs,
		Members:            members,
		Settings:           settings,
		Background:         channel.Background,
		Status:             channel.Status,
		UnreadMessageCount: unreadCount,
		CreatedAt:          channel.CreatedAt,
		UpdatedAt:          channel.UpdatedAt,
	}
}

func FromChannels(channels []model.Channel, unreadCounts []*int64) ([]ChannelResponse, error) {
	if len(channels) != len(unreadCounts) {
		return nil, errors.New("channels and unread message have different lengths")
	}

	responses := make([]ChannelResponse, len(channels))
	for i, ch := range channels {
		responses[i] = FromChannel(&ch, unreadCounts[i])
	}
	return responses, nil
}
