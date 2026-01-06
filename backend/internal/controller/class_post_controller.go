package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
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
// Content-Type: application/json
//
// Workflow:
// 1. Upload files first: POST /api/classrooms/posts/upload-attachments
// 2. Get file URLs from step 1
// 3. Create post with URLs in attachments array
//
// Request body:
// {
//   "title": "Post title",
//   "content": "Post content",
//   "attachments": [
//     {
//       "file_name": "document.pdf",
//       "file_url": "https://cloudinary.com/...",
//       "file_size": 1024000,
//       "mime_type": "application/pdf"
//     }
//   ]
// }

func (c *ClassPostController) CreatePost(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.CreateClassPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	post, err := c.postService.CreatePost(req, classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Post created successfully", dto.FromClassPost(post))
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

	posts, total, err := c.postService.GetPostsByClassroom(classroomID, user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"posts":     dto.FromClassPosts(posts),
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

	post, err := c.postService.GetPostByID(postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post retrieved successfully", dto.FromClassPost(post))
}

// UpdatePost updates a post
// PUT /api/classrooms/posts/:post_id
// Content-Type: application/json
func (c *ClassPostController) UpdatePost(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.UpdateClassPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	post, err := c.postService.UpdatePost(req, postID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post updated successfully", dto.FromClassPost(post))
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

// UploadAttachments uploads files to Cloudinary and returns URLs
// POST /api/classrooms/posts/upload-attachments
// This should be called BEFORE creating/updating a post
func (c *ClassPostController) UploadAttachments(ctx *gin.Context) {
	_, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Get uploaded files
	form, err := ctx.MultipartForm()
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "No files uploaded", "NO_FILES")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		dto.SendError(ctx, http.StatusBadRequest, "No files uploaded", "NO_FILES")
		return
	}

	// Upload files to Cloudinary
	uploadedFiles := make([]gin.H, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		// Upload to Cloudinary
		result, err := cloudinary.Upload(file)
		file.Close()

		if err != nil {
			// Skip failed uploads
			continue
		}

		uploadedFiles = append(uploadedFiles, gin.H{
			"file_name": fileHeader.Filename,
			"file_url":  result.SecureURL,
			"file_size": fileHeader.Size,
			"mime_type": fileHeader.Header.Get("Content-Type"),
			"public_id": result.PublicID,
		})
	}

	if len(uploadedFiles) == 0 {
		dto.SendError(ctx, http.StatusInternalServerError, "All file uploads failed", "UPLOAD_FAILED")
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Files uploaded successfully", gin.H{
		"attachments": uploadedFiles,
	})
}

// DeleteAttachment deletes a file from Cloudinary by public_id
// DELETE /api/classrooms/posts/attachments/:public_id
func (c *ClassPostController) DeleteAttachment(ctx *gin.Context) {
	_, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	publicID := ctx.Param("public_id")
	if publicID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Public ID is required", apperror.ErrBadRequest.Code)
		return
	}

	// Delete from Cloudinary
	result, err := cloudinary.Delete(publicID)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to delete file", "DELETE_FAILED")
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "File deleted successfully", gin.H{
		"public_id": publicID,
		"result":    result.Result,
	})
}

// AddAttachment adds an attachment to a post
// POST /api/classrooms/posts/:post_id/attachments
func (c *ClassPostController) AddAttachment(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req struct {
		FileName string `json:"file_name" binding:"required"`
		FileURL  string `json:"file_url" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
		MimeType string `json:"mime_type" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	attachment := model.Attachment{
		FileName: req.FileName,
		FileURL:  req.FileURL,
		FileSize: req.FileSize,
		MimeType: req.MimeType,
	}

	post, err := c.postService.AddAttachment(postID, user.ID, attachment)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Attachment added successfully", dto.FromClassPost(post))
}

// RemoveAttachment removes an attachment from a post
// DELETE /api/classrooms/posts/:post_id/attachments
func (c *ClassPostController) RemoveAttachment(ctx *gin.Context) {
	postID := ctx.Param("post_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req struct {
		FileURL string `json:"file_url" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	post, err := c.postService.RemoveAttachment(postID, user.ID, req.FileURL)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Attachment removed successfully", dto.FromClassPost(post))
}
