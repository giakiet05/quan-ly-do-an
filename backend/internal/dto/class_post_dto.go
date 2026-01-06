package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// ===== REQUEST DTOs =====

type CreateClassPostRequest struct {
	Title       string             `json:"title" binding:"required"`
	Content     string             `json:"content" binding:"required"`
	Attachments []AttachmentUpload `json:"attachments,omitempty"`
}

type UpdateClassPostRequest struct {
	Title              *string            `json:"title"`
	Content            *string            `json:"content"`
	AttachmentsToAdd   []AttachmentUpload `json:"attachments_to_add,omitempty"`
	AttachmentsToRemove []string          `json:"attachments_to_remove,omitempty"` // file URLs to remove
}

type AttachmentUpload struct {
	FileName string `json:"file_name" binding:"required"`
	FileURL  string `json:"file_url" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required"`
	MimeType string `json:"mime_type" binding:"required"`
}

// ===== RESPONSE DTOs =====

type ClassPostResponse struct {
	ID          string                 `json:"id"`
	ClassroomID string                 `json:"classroom_id"`
	Author      model.UserInfoResponse `json:"author"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Attachments []model.Attachment     `json:"attachments,omitempty"`
	IsPinned    bool                   `json:"is_pinned"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func FromClassPost(post *model.ClassPost) ClassPostResponse {
	if post == nil {
		return ClassPostResponse{}
	}
	return ClassPostResponse{
		ID:          post.ID.Hex(),
		ClassroomID: post.ClassroomID.Hex(),
		Author: model.UserInfoResponse{
			UserID:   post.Author.ID.Hex(),
			FullName: post.Author.FullName,
			Avatar:   post.Author.Avatar,
		},
		Title:       post.Title,
		Content:     post.Content,
		Attachments: post.Attachments,
		IsPinned:    post.IsPinned,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

func FromClassPosts(posts []model.ClassPost) []ClassPostResponse {
	responses := make([]ClassPostResponse, 0, len(posts))
	for _, post := range posts {
		responses = append(responses, FromClassPost(&post))
	}
	return responses
}
