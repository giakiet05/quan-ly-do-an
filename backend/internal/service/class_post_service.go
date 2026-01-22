package service

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/cloudinary"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassPostService interface {
	CreatePost(req dto.CreateClassPostRequest, classroomID string, authorID string) (*dto.ClassPostResponse, error)
	GetPostByID(postID string) (*dto.ClassPostResponse, error)
	GetPostsByClassroom(classroomID string, userID string, page, pageSize int) ([]dto.ClassPostResponse, int64, error)
	UpdatePost(req dto.UpdateClassPostRequest, postID string, userID string) (*dto.ClassPostResponse, error)
	DeletePost(postID string, userID string) error
	TogglePinPost(postID string, userID string, isPinned bool) error
}

type classPostService struct {
	postRepo      repo.ClassPostRepo
	classroomRepo repo.ClassroomRepo
	userRepo      repo.UserRepo
	eventBus      *bus.EventBus
}

func NewClassPostService(
	postRepo repo.ClassPostRepo,
	classroomRepo repo.ClassroomRepo,
	userRepo repo.UserRepo,
	eventBus *bus.EventBus,
) ClassPostService {
	return &classPostService{
		postRepo:      postRepo,
		classroomRepo: classroomRepo,
		userRepo:      userRepo,
		eventBus:      eventBus,
	}
}

// Helper: Publish ClassPostCreatedEvent
func (s *classPostService) publishClassPostCreatedEvent(ctx context.Context, classroomID string, post *model.ClassPost, author *model.User) {
	// Get classroom to get student IDs
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return
	}

	// Don't send notification if classroom has no students
	if len(classroom.StudentIDs) == 0 {
		return
	}

	s.eventBus.Publish(bus.ClassPostCreatedEvent{
		ClassroomID: classroomID,
		PostID:      post.ID.Hex(),
		PostTitle:   post.Title,
		AuthorID:    author.ID.Hex(),
		AuthorName:  author.FullName,
	})
}

// Helper: Publish ClassPostUpdatedEvent
func (s *classPostService) publishClassPostUpdatedEvent(ctx context.Context, classroomID string, post *model.ClassPost, author *model.User) {
	// Get classroom to get student IDs
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return
	}

	// Don't send notification if classroom has no students
	if len(classroom.StudentIDs) == 0 {
		return
	}

	s.eventBus.Publish(bus.ClassPostUpdatedEvent{
		ClassroomID: classroomID,
		PostID:      post.ID.Hex(),
		PostTitle:   post.Title,
		AuthorID:    author.ID.Hex(),
		AuthorName:  author.FullName,
	})
}

// Helper: Populate author for a single post
func (s *classPostService) populatePostAuthor(ctx context.Context, post *model.ClassPost) (*model.User, error) {
	return s.userRepo.GetByID(ctx, post.AuthorID.Hex())
}

// Helper: Batch populate authors for multiple posts
func (s *classPostService) populatePostsAuthors(ctx context.Context, posts []model.ClassPost) (map[string]*model.User, error) {
	// Collect unique author IDs
	authorIDsMap := make(map[string]bool)
	for _, post := range posts {
		authorIDsMap[post.AuthorID.Hex()] = true
	}

	// Convert to slice
	authorIDs := make([]string, 0, len(authorIDsMap))
	for id := range authorIDsMap {
		authorIDs = append(authorIDs, id)
	}

	// Batch query users
	users, err := s.userRepo.GetByIDs(ctx, authorIDs)
	if err != nil {
		return nil, err
	}

	// Create map
	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.ID.Hex()] = user
	}

	return userMap, nil
}

func (s *classPostService) CreatePost(req dto.CreateClassPostRequest, classroomID string, authorID string) (*dto.ClassPostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer of the classroom
	isLecturer, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, authorID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	// Get author info for response
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	classroomOID, _ := primitive.ObjectIDFromHex(classroomID)
	authorOID, _ := primitive.ObjectIDFromHex(authorID)

	// Convert attachments
	attachments := make([]model.Attachment, 0, len(req.Attachments))
	for _, att := range req.Attachments {
		attachments = append(attachments, model.Attachment{
			FileName: att.FileName,
			FileURL:  att.FileURL,
			PublicID: att.PublicID,
			FileSize: att.FileSize,
			MimeType: att.MimeType,
		})
	}

	post := &model.ClassPost{
		ClassroomID: classroomOID,
		AuthorID:    authorOID,
		Title:       req.Title,
		Content:     req.Content,
		Attachments: attachments,
		IsPinned:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdPost, err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	// Publish event for notification
	s.publishClassPostCreatedEvent(ctx, classroomID, createdPost, author)

	// Return DTO with author info
	response := dto.FromClassPostWithAuthor(createdPost, author)
	return &response, nil
}

func (s *classPostService) GetPostByID(postID string) (*dto.ClassPostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Populate author
	author, err := s.populatePostAuthor(ctx, post)
	if err != nil {
		return nil, err
	}

	response := dto.FromClassPostWithAuthor(post, author)
	return &response, nil
}

func (s *classPostService) GetPostsByClassroom(classroomID string, userID string, page, pageSize int) ([]dto.ClassPostResponse, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is in classroom (student or lecturer)
	isLecturer, _ := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, userID)
	isStudent, _ := s.classroomRepo.IsStudentInClassroom(ctx, classroomID, userID)

	if !isLecturer && !isStudent {
		return nil, 0, apperror.ErrForbidden
	}

	posts, total, err := s.postRepo.GetByClassroom(ctx, classroomID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Batch populate authors
	authorsMap, err := s.populatePostsAuthors(ctx, posts)
	if err != nil {
		return nil, 0, err
	}

	// Convert to DTOs
	responses := dto.FromClassPostsWithAuthors(posts, authorsMap)
	return responses, total, nil
}

func (s *classPostService) UpdatePost(req dto.UpdateClassPostRequest, postID string, userID string) (*dto.ClassPostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check if user is the author or lecturer
	isLecturer, _ := s.classroomRepo.IsLecturerOrCoLecturer(ctx, post.ClassroomID.Hex(), userID)
	isAuthor := post.AuthorID.Hex() == userID

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
				PublicID: att.PublicID,
				FileSize: att.FileSize,
				MimeType: att.MimeType,
			})
		}
	}

	post.UpdatedAt = time.Now()

	updatedPost, err := s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	// Populate author for response
	author, err := s.populatePostAuthor(ctx, updatedPost)
	if err != nil {
		return nil, err
	}

	// Publish event for notification
	s.publishClassPostUpdatedEvent(ctx, post.ClassroomID.Hex(), updatedPost, author)

	response := dto.FromClassPostWithAuthor(updatedPost, author)
	return &response, nil
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
	isLecturer, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, post.ClassroomID.Hex(), userID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	// Delete post from DB first
	err = s.postRepo.Delete(ctx, postID)
	if err != nil {
		return err
	}

	// Delete attachments from Cloudinary (best effort, ignore errors)
	for _, att := range post.Attachments {
		if att.PublicID != "" {
			_, _ = cloudinary.Delete(att.PublicID)
		}
	}

	return nil
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
	isLecturer, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, post.ClassroomID.Hex(), userID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	return s.postRepo.TogglePin(ctx, postID, isPinned)
}

