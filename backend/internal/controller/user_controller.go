package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/cloudinary"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// UserController handles requests related to user management.
type UserController struct {
	service service.UserService
}

// NewUserController creates a new UserController.
func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

// GetUsers retrieves a paginated list of users with optional search filters.
func (c *UserController) GetUsers(ctx *gin.Context) {
	var query dto.GetUsersQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	response, err := c.service.GetUsers(&query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, "Users retrieved successfully", response)
}

// GetMyProfile retrieves the profile of the currently authenticated user.
func (c *UserController) GetMyProfile(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	user, err := c.service.GetUserByID(authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Profile retrieved successfully", user)
}

// UploadAvatar handles avatar image uploads.
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid form data", "INVALID_FORM")
		return
	}

	images, err := cloudinary.UploadImages(form.File["avatar"])
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload image", "UPLOAD_FAILED")
		return
	}

	updatedUser, err := c.service.UpdateAvatar(authUser.(auth.AuthUser).ID, images[0].URL, images[0].PublicID)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to update avatar", "DB_UPDATE_FAILED")
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Avatar updated successfully", updatedUser)
}

// DeleteAvatar removes the user's avatar.
func (c *UserController) DeleteAvatar(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	updatedUser, err := c.service.DeleteAvatar(authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Avatar deleted successfully", updatedUser)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.service.ChangePassword(authUser.(auth.AuthUser).ID, req.OldPassword, req.NewPassword)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Password changed successfully", nil)
}

// UpdateProfile updates user profile (full_name, student_code)
// PUT /api/users/me
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, err := c.service.UpdateProfile(authUser.(auth.AuthUser).ID, req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Profile updated successfully", user)
}

// --- Admin-only actions ---

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	err := c.service.DeleteUser(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User deleted successfully", gin.H{"id": userID})
}
