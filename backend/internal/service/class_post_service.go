package service

import (
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassPostService interface {
	CreatePost(req dto.CreateClassPostRequest, classroomID string, authorID string) (*model.ClassPost, error)
	GetPostByID(postID string) (*model.ClassPost, error)
	GetPostsByClassroom(classroomID string, userID string, page, pageSize int) ([]model.ClassPost, int64, error)
	UpdatePost(req dto.UpdateClassPostRequest, postID string, userID string) (*model.ClassPost, error)
	DeletePost(postID string, userID string) error
	TogglePinPost(postID string, userID string, isPinned bool) error
	AddAttachment(postID string, userID string, attachment model.Attachment) (*model.ClassPost, error)
	RemoveAttachment(postID string, userID string, fileURL string) (*model.ClassPost, error)
}

type classPostService struct {
	postRepo      repo.ClassPostRepo
	classroomRepo repo.ClassroomRepo
	userRepo      repo.UserRepo
}

func NewClassPostService(
	postRepo repo.ClassPostRepo,
	classroomRepo repo.ClassroomRepo,
	userRepo repo.UserRepo,
) ClassPostService {
	return &classPostService{
		postRepo:      postRepo,
		classroomRepo: classroomRepo,
		userRepo:      userRepo,
	}
}

func (s *classPostService) CreatePost(req dto.CreateClassPostRequest, classroomID string, authorID string) (*model.ClassPost, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer of the classroom
	isLecturer, err := s.classroomRepo.IsLecturer(ctx, classroomID, authorID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	// Get author info
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	classroomOID, _ := primitive.ObjectIDFromHex(classroomID)

	// Convert attachments
	attachments := make([]model.Attachment, 0, len(req.Attachments))
	for _, att := range req.Attachments {
		attachments = append(attachments, model.Attachment{
			FileName: att.FileName,
			FileURL:  att.FileURL,
			FileSize: att.FileSize,
			MimeType: att.MimeType,
		})
	}

	post := &model.ClassPost{
		ClassroomID: classroomOID,
		Author: model.UserInfo{
			ID:       author.ID,
			FullName: author.FullName,
			Avatar:   author.Avatar,
		},
		Title:       req.Title,
		Content:     req.Content,
		Attachments: attachments,
		IsPinned:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.postRepo.Create(ctx, post)
}

func (s *classPostService) GetPostByID(postID string) (*model.ClassPost, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.postRepo.GetByID(ctx, postID)
}

func (s *classPostService) GetPostsByClassroom(classroomID string, userID string, page, pageSize int) ([]model.ClassPost, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is in classroom (student or lecturer)
	isLecturer, _ := s.classroomRepo.IsLecturer(ctx, classroomID, userID)
	isStudent, _ := s.classroomRepo.IsStudentInClassroom(ctx, classroomID, userID)

	if !isLecturer && !isStudent {
		return nil, 0, apperror.ErrForbidden
	}

	return s.postRepo.GetByClassroom(ctx, classroomID, page, pageSize)
}

func (s *classPostService) UpdatePost(req dto.UpdateClassPostRequest, postID string, userID string) (*model.ClassPost, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check if user is the author or lecturer
	isLecturer, _ := s.classroomRepo.IsLecturer(ctx, post.ClassroomID.Hex(), userID)
	isAuthor := post.Author.ID.Hex() == userID

	if !isLecturer && !isAuthor {
		return nil, apperror.ErrForbidden
	}

	// Update fields
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}

	// Remove attachments
	if len(req.AttachmentsToRemove) > 0 {
		newAttachments := make([]model.Attachment, 0)
		for _, att := range post.Attachments {
			shouldRemove := false
			for _, urlToRemove := range req.AttachmentsToRemove {
				if att.FileURL == urlToRemove {
					shouldRemove = true
					break
				}
			}
			if !shouldRemove {
				newAttachments = append(newAttachments, att)
			}
		}
		post.Attachments = newAttachments
	}

	// Add new attachments
	if len(req.AttachmentsToAdd) > 0 {
		for _, att := range req.AttachmentsToAdd {
			post.Attachments = append(post.Attachments, model.Attachment{
				FileName: att.FileName,
				FileURL:  att.FileURL,
				FileSize: att.FileSize,
				MimeType: att.MimeType,
			})
		}
	}

	post.UpdatedAt = time.Now()

	return s.postRepo.Update(ctx, post)
}

func (s *classPostService) DeletePost(postID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Only lecturer can delete
	isLecturer, err := s.classroomRepo.IsLecturer(ctx, post.ClassroomID.Hex(), userID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	return s.postRepo.Delete(ctx, postID)
}

func (s *classPostService) TogglePinPost(postID string, userID string, isPinned bool) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Only lecturer can pin/unpin
	isLecturer, err := s.classroomRepo.IsLecturer(ctx, post.ClassroomID.Hex(), userID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	return s.postRepo.TogglePin(ctx, postID, isPinned)
}

func (s *classPostService) AddAttachment(postID string, userID string, attachment model.Attachment) (*model.ClassPost, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check if user is the author or lecturer
	isLecturer, _ := s.classroomRepo.IsLecturer(ctx, post.ClassroomID.Hex(), userID)
	isAuthor := post.Author.ID.Hex() == userID

	if !isLecturer && !isAuthor {
		return nil, apperror.ErrForbidden
	}

	// Add attachment to array
	post.Attachments = append(post.Attachments, attachment)
	post.UpdatedAt = time.Now()

	return s.postRepo.Update(ctx, post)
}

func (s *classPostService) RemoveAttachment(postID string, userID string, fileURL string) (*model.ClassPost, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check if user is the author or lecturer
	isLecturer, _ := s.classroomRepo.IsLecturer(ctx, post.ClassroomID.Hex(), userID)
	isAuthor := post.Author.ID.Hex() == userID

	if !isLecturer && !isAuthor {
		return nil, apperror.ErrForbidden
	}

	// Filter out the attachment
	newAttachments := make([]model.Attachment, 0)
	for _, att := range post.Attachments {
		if att.FileURL != fileURL {
			newAttachments = append(newAttachments, att)
		}
	}

	post.Attachments = newAttachments
	post.UpdatedAt = time.Now()

	return s.postRepo.Update(ctx, post)
}
