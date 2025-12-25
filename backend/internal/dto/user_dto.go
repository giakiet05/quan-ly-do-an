package dto

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// --- Request DTOs ---

// New Registration Flow (Verify Email First)

// GetUsersQuery contains query parameters for searching and paginating users
type GetUsersQuery struct {
	Username string `form:"username"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// --- Response DTOs ---

// UserResponse is the main user object returned in API responses.
type UserResponse struct {
	ID       string             `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email,omitempty"`
	Provider model.AuthProvider `json:"provider"`
	Avatar   *model.Image       `json:"avatar,omitempty"`
}

func FromUser(u *model.User) UserResponse {
	if u == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Provider: u.Provider,
		Avatar:   u.Avatar,
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
