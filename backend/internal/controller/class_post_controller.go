package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/cloudinary"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ClassPostController struct {
	postService service.ClassPostService
}

func NewClassPostController(postService service.ClassPostService) *ClassPostController {
	return &ClassPostController{
		postService: postService,
	}
}

// CreatePost creates a new class post (lecturer only)
// POST /api/classrooms/:id/posts
// Content-Type: multipart/form-data
//
// Form fields:
// - title: string (required)
// - content: string (required)
// - files: file[] (optional, multiple files)

func (c *ClassPostController) CreatePost(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse form data
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")

	if title == "" || content == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Title and content are required", apperror.ErrBadRequest.Code)
		return
	}

	// Get uploaded files
	form, err := ctx.MultipartForm()
	var attachments []dto.AttachmentUpload
	uploadedPublicIDs := []string{}

	if err == nil && form != nil {
		files := form.File["files"]

		// Upload files to Cloudinary
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to open file", "FILE_OPEN_FAILED")
				return
			}

			result, err := cloudinary.Upload(file)
			file.Close()

			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload file", "FILE_UPLOAD_FAILED")
				return
			}

			uploadedPublicIDs = append(uploadedPublicIDs, result.PublicID)
			attachments = append(attachments, dto.AttachmentUpload{
				FileName: fileHeader.Filename,
				FileURL:  result.SecureURL,
				PublicID: result.PublicID,
				FileSize: fileHeader.Size,
				MimeType: fileHeader.Header.Get("Content-Type"),
			})
		}
	}

	// Create post request
	req := dto.CreateClassPostRequest{
		Title:       title,
		Content:     content,
		Attachments: attachments,
	}

	response, err := c.postService.CreatePost(req, classroomID, user.ID)
	if err != nil {
		// Rollback: delete uploaded files
		for _, pid := range uploadedPublicIDs {
			_, _ = cloudinary.Delete(pid)
		}
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Post created successfully", response)
}

// GetPosts gets all posts in a classroom
// GET /api/classrooms/:id/posts
func (c *ClassPostController) GetPosts(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse pagination
	var query struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	responses, total, err := c.postService.GetPostsByClassroom(classroomID, user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"posts":     responses,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	}

	dto.SendSuccess(ctx, http.StatusOK, "Posts retrieved successfully", data)
}

// GetPost gets a single post by ID
// GET /api/classrooms/posts/:post_id
func (c *ClassPostController) GetPost(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	response, err := c.postService.GetPostByID(postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post retrieved successfully", response)
}

// UpdatePost updates a post
// PUT /api/classrooms/posts/:post_id
// Content-Type: multipart/form-data
//
// Form fields:
// - title: string (optional)
// - content: string (optional)
// - files: file[] (optional, files to add)
// - files_to_remove: string (optional, JSON array of URLs to remove)
func (c *ClassPostController) UpdatePost(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse form data
	var req dto.UpdateClassPostRequest

	if title := ctx.PostForm("title"); title != "" {
		req.Title = &title
	}
	if content := ctx.PostForm("content"); content != "" {
		req.Content = &content
	}

	// Parse files_to_remove (JSON array string)
	if removeStr := ctx.PostForm("files_to_remove"); removeStr != "" {
		var filesToRemove []string
		if err := json.Unmarshal([]byte(removeStr), &filesToRemove); err == nil {
			req.AttachmentsToRemove = filesToRemove
		}
	}

	// Get uploaded files
	form, err := ctx.MultipartForm()
	uploadedPublicIDs := []string{}

	if err == nil && form != nil {
		files := form.File["files"]

		// Upload new files to Cloudinary
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to open file", "FILE_OPEN_FAILED")
				return
			}

			result, err := cloudinary.Upload(file)
			file.Close()

			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload file", "FILE_UPLOAD_FAILED")
				return
			}

			uploadedPublicIDs = append(uploadedPublicIDs, result.PublicID)
			req.AttachmentsToAdd = append(req.AttachmentsToAdd, dto.AttachmentUpload{
				FileName: fileHeader.Filename,
				FileURL:  result.SecureURL,
				PublicID: result.PublicID,
				FileSize: fileHeader.Size,
				MimeType: fileHeader.Header.Get("Content-Type"),
			})
		}
	}

	response, err := c.postService.UpdatePost(req, postID, user.ID)
	if err != nil {
		// Rollback: delete newly uploaded files
		for _, pid := range uploadedPublicIDs {
			cloudinary.Delete(pid)
		}
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	// Delete removed files from Cloudinary (best effort)
	for _, url := range req.AttachmentsToRemove {
		publicID := extractPublicIDFromURL(url)
		if publicID != "" {
			_, _ = cloudinary.Delete(publicID)
		}
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post updated successfully", response)
}

// DeletePost deletes a post (lecturer only)
// DELETE /api/classrooms/posts/:post_id
func (c *ClassPostController) DeletePost(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	err := c.postService.DeletePost(postID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post deleted successfully", gin.H{"post_id": postID})
}

// TogglePinPost pins/unpins a post (lecturer only)
// PATCH /api/classrooms/posts/:post_id/pin
func (c *ClassPostController) TogglePinPost(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req struct {
		IsPinned bool `json:"is_pinned"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.postService.TogglePinPost(postID, user.ID, req.IsPinned)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post pin status updated", gin.H{"is_pinned": req.IsPinned})
}


// Helper function to extract public_id from Cloudinary URL
// URL format: https://res.cloudinary.com/xxx/image/upload/v123456/folder/filename.jpg
// Returns: folder/filename
func extractPublicIDFromURL(url string) string {
	// Find "/upload/" in URL
	parts := strings.Split(url, "/upload/")
	if len(parts) < 2 {
		return ""
	}

	// Get path after /upload/v123456/
	path := parts[1]
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) < 2 {
		return ""
	}

	// Remove file extension
	publicID := pathParts[1]
	lastDot := strings.LastIndex(publicID, ".")
	if lastDot > 0 {
		publicID = publicID[:lastDot]
	}

	return publicID
}
