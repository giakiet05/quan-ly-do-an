package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

type CreateMessageRequest struct {
	ChannelID string            `json:"channel_id"`
	SenderID  string            `json:"sender_id"`
	Type      model.MessageType `json:"type"`
	Content   string            `json:"content"`
}

type GetMessageFilterQuery struct {
	ChannelID     string  `form:"channel_id"`
	SenderID      *string `form:"sender_id,omitempty"`
	SearchContent *string `form:"search_content,omitempty"`
	IsRead        *bool   `form:"is_read,omitempty"`
	IsSend        *bool   `form:"is_send,omitempty"`
	IsMedia       *bool   `form:"is_media,omitempty"`
	Page          int     `form:"page"`
	PageSize      int     `form:"page_size"`
}

type MessageResponse struct {
	ID        string            `json:"id"`
	ChannelID string            `json:"channel_id"`
	SenderID  string            `json:"sender_id"`
	Type      model.MessageType `json:"type"`
	Content   string            `json:"content"`
	ReadBy    []string          `json:"read_by"`
	CreatedAt time.Time         `json:"created_at"`
}

func FromMessage(message *model.Message) MessageResponse {
	if message == nil {
		return MessageResponse{}
	}

	var senderID string
	if message.SenderID != nil {
		senderID = message.SenderID.Hex()
	}

	readByString := make([]string, 0, len(message.ReadBy))
	for _, readBy := range message.ReadBy {
		readByString = append(readByString, readBy.Hex())
	}

	return MessageResponse{
		ID:        message.ID.Hex(),
		ChannelID: message.ChannelID.Hex(),
		SenderID:  senderID,
		Type:      message.Type,
		Content:   message.Content,
		ReadBy:    readByString,
		CreatedAt: message.CreatedAt,
	}
}

func FromMessages(messages []model.Message) []MessageResponse {
	messageResponses := make([]MessageResponse, 0, len(messages))
	for _, msg := range messages {
		messageResponses = append(messageResponses, FromMessage(&msg))
	}
	return messageResponses
}
