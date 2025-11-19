package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type CreateMessageRequest struct {
	ChannelID string            `json:"channel_id"`
	SenderID  string            `json:"sender_id"`
	Type      model.MessageType `json:"type"`
	Content   string            `json:"content"`
}

type GetMessageFilterQuery struct {
	ChannelID     string  `json:"channel_id"`
	SenderID      *string `json:"sender_id,omitempty"`
	SearchContent *string `json:"search_content,omitempty"`
	IsRead        *bool   `json:"is_read,omitempty"`
	IsSend        *bool   `json:"is_send,omitempty"`
	IsMedia       *bool   `json:"is_media,omitempty"`
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
}

type MessageResponse struct {
	ID             string            `json:"id"`
	ChannelID      string            `json:"channel_id"`
	SenderID       string            `json:"sender_id"`
	SenderUsername string            `json:"sender_username"`
	Type           model.MessageType `json:"type"`
	Content        string            `json:"content"`
	ReadBy         []string          `json:"read_by"`
	CreatedAt      time.Time         `json:"created_at"`
}

func FromMessage(message *model.Message) *MessageResponse {
	var senderID string
	if message.SenderID != nil {
		senderID = message.SenderID.Hex()
	}

	readByString := make([]string, 0, len(message.ReadBy))
	for _, readBy := range message.ReadBy {
		readByString = append(readByString, readBy.Hex())
	}

	return &MessageResponse{
		ID:             message.ID.Hex(),
		ChannelID:      message.ChannelID.Hex(),
		SenderID:       senderID,
		SenderUsername: message.SenderUsername,
		Type:           message.Type,
		Content:        message.Content,
		ReadBy:         readByString,
		CreatedAt:      message.CreatedAt,
	}
}

func FromMessages(messages []model.Message) []MessageResponse {
	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, *FromMessage(&msg))
	}
	return messageResponses
}
