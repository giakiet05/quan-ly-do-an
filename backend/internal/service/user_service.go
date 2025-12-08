package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/cloudinary"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to user management.
type UserService interface {
	UpdateAvatar(userID string, imageURL string, publicID string) (dto.UserResponse, error)
	DeleteAvatar(userID string) (dto.UserResponse, error)
	DeleteUser(id string) error
	ChangePassword(userID, oldPassword, newPassword string) error

	GetUserByID(id string) (dto.UserResponse, error)
	GetUserByUsername(username string, requesterID string) (dto.UserResponse, error)
	GetUserByEmail(email string) (dto.UserResponse, error)
	GetUsers(query *dto.GetUsersQuery) (*dto.PaginatedUsersResponse, error)

	CheckUsernameAvailability(username string) (bool, error)
}

type userService struct {
	userRepo    repo.UserRepo
	eventBus    *bus.EventBus
	redisClient *redis.Client
}

func NewUserService(userRepo repo.UserRepo, bus *bus.EventBus, redisClient *redis.Client) UserService {
	return &userService{
		userRepo:    userRepo,
		eventBus:    bus,
		redisClient: redisClient,
	}
}

func (s *userService) UpdateAvatar(userID string, imageURL string, publicID string) (dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	var oldPublicID string
	newImage := &model.Image{URL: imageURL, PublicID: publicID}

	//Lay publicid cu de xoa
	if user.Avatar != nil {
		oldPublicID = user.Avatar.PublicID
	}
	user.Avatar = newImage

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if oldPublicID != "" {
		go cloudinary.Delete(oldPublicID)
	}

	s.eventBus.Publish(bus.UserChangeAvatarEventType{UserID: userID, NewAvatar: imageURL})

	return dto.FromUser(updatedUser), nil
}

func (s *userService) DeleteAvatar(userID string) (dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	var oldPublicID string
	if user.Avatar != nil {
		oldPublicID = user.Avatar.PublicID
	}
	user.Avatar = nil

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if oldPublicID != "" {
		go cloudinary.Delete(oldPublicID)
	}

	s.eventBus.Publish(bus.UserChangeAvatarEventType{UserID: userID, NewAvatar: ""})

	return dto.FromUser(updatedUser), nil
}

func (s *userService) DeleteUser(id string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if auth.TokenSvc != nil {
		if err := auth.TokenSvc.InvalidateAllUserTokens(ctx, id); err != nil {
			fmt.Printf("Failed to invalidate tokens for user %s: %v\n", id, err)
		}
	}

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) ChangePassword(userID, oldPassword, newPassword string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	if user.Provider != model.ProviderLocal || user.Password == "" {
		return apperror.ErrLoginMethodMismatch
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return apperror.ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) GetUserByID(id string) (dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.UserResponse{}, apperror.ErrUserNotFound
		}
		return dto.UserResponse{}, err
	}
	return dto.FromUser(user), nil
}

func (s *userService) GetUserByUsername(username string, requesterID string) (dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.UserResponse{}, apperror.ErrUserNotFound
		}
		return dto.UserResponse{}, err
	}

	return dto.FromUser(user), nil
}

func (s *userService) GetUserByEmail(email string) (dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.UserResponse{}, apperror.ErrUserNotFound
		}
		return dto.UserResponse{}, err
	}
	return dto.FromUser(user), nil
}

func (s *userService) GetUsers(query *dto.GetUsersQuery) (*dto.PaginatedUsersResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Build filter - exclude deleted and banned users for regular users
	filter := repo.Filter{
		"deleted_at": bson.M{"$exists": false},
		"is_banned":  false, // Only show non-banned users
	}

	// Add username search if provided
	if query.Username != "" {
		filter["username"] = bson.M{"$regex": primitive.Regex{Pattern: query.Username, Options: "i"}}
	}

	// Pagination
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	findOptions := &repo.FindOptions{
		Skip:  int64((page - 1) * pageSize),
		Limit: int64(pageSize),
		Sort:  map[string]int{"created_at": -1},
	}

	users, total, err := s.userRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	userResponses := dto.FromUsers(users)

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *userService) CheckUsernameAvailability(username string) (bool, error) {
	// Try cache first
	if s.redisClient != nil {
		ctx, cancel := util.NewDefaultRedisContext()
		defer cancel()

		cacheKey := fmt.Sprintf("username_exists:%s", username)
		cached, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache hit - "false" means available, "true" means taken
			return cached == "false", nil
		}
	}

	// Cache miss - query database
	dbCtx, cancel := util.NewDefaultDBContext()
	defer cancel()

	_, err := s.userRepo.GetByUsername(dbCtx, username)
	exists := !errors.Is(err, mongo.ErrNoDocuments)

	// Cache the result (5 minutes TTL)
	if s.redisClient != nil {
		ctx, cancel := util.NewDefaultRedisContext()
		defer cancel()

		cacheKey := fmt.Sprintf("username_exists:%s", username)
		value := "false"
		if exists {
			value = "true"
		}
		// Ignore cache write errors, not critical
		_ = s.redisClient.Set(ctx, cacheKey, value, 5*time.Minute).Err()
	}

	return !exists, nil
}
