package dto

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// --- Request DTOs ---

// New Registration Flow (Verify Email First)

// GetUsersQuery contains query parameters for searching and paginating users
type GetUsersQuery struct {
	FullName string `form:"full_name"`
	Email    string `form:"email"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	FullName    *string `json:"full_name"`
	StudentCode *string `json:"student_code"`
}

// --- Response DTOs ---

// UserResponse is the main user object returned in API responses.
type UserResponse struct {
	ID          string             `json:"id"`
	Email       string             `json:"email,omitempty"`
	FullName    string             `json:"full_name"`
	StudentCode *string            `json:"student_code,omitempty"`
	Provider    model.AuthProvider `json:"provider"`
	Avatar      *model.Image       `json:"avatar,omitempty"`
	CreatedAt   string             `json:"created_at,omitempty"`
}

func FromUser(u *model.User) UserResponse {
	if u == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:          u.ID.Hex(),
		Email:       u.Email,
		FullName:    u.FullName,
		StudentCode: u.StudentCode,
		Provider:    u.AuthProvider,
		Avatar:      u.Avatar,
		CreatedAt:   u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func FromUsers(users []*model.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		if u == nil {
			continue // Skip nil users
		}
		userResponse := FromUser(u)
		//userResponse.Email = "" // Hide email in list view
		responses = append(responses, userResponse)
	}
	return responses
}

type UserInfoResponse struct {
	UserID      string       `json:"user_id"`
	FullName    string       `json:"full_name"`
	Email       string       `json:"email"`
	Avatar      *model.Image `json:"avatar,omitempty"`
	StudentCode *string      `json:"student_code,omitempty"`
}

func ToUserInfoResponse(u *model.User) UserInfoResponse {
	if u == nil {
		return UserInfoResponse{}
	}
	return UserInfoResponse{
		UserID:      u.ID.Hex(),
		FullName:    u.FullName,
		Email:       u.Email,
		Avatar:      u.Avatar,
		StudentCode: u.StudentCode,
	}
}

func ToUserInfoResponses(users []*model.User) []UserInfoResponse {
	responses := make([]UserInfoResponse, 0, len(users))
	for _, u := range users {
		if u == nil {
			continue
		}
		responses = append(responses, ToUserInfoResponse(u))
	}
	return responses
}
