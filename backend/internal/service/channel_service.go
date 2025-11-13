package service

import (
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelService interface {
	Start()
	CreateChannel(req *dto.CreateChannelRequest, requesterID string) (*model.Channel, error)
	GetChannelByID(channelID string) (*model.Channel, error)
	GetChannelsByUserID(userID string, requesterID string, page int, pageSize int) (*dto.PaginatedChannelsResponse, error)
	GetChannelByBothUserID(user1ID string, user2ID string, requesterID string) (*model.Channel, error)
	UpdateChannel(req *dto.UpdateChannelRequest, requesterID string) (*model.Channel, error)
	DeleteChannel(channelID string, userID string) error
}

type channelService struct {
	channelRepository repo.ChannelRepo
	eventBus          *bus.EventBus
}

func NewChannelService(channelRepo repo.ChannelRepo, bus *bus.EventBus) ChannelService {
	return &channelService{
		channelRepository: channelRepo,
		eventBus:          bus,
	}
}

func (s *channelService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.eventBus.Subscribe(bus.TopicUserChangeAvatar, eventChannel)

	log.Println("ChannelService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

func (s *channelService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicUserChangeAvatar:
			s.handleNewAvatar(event)
		default:
			log.Println("Unhandled event topic:", event.Topic())
		}
	}
}

func (s *channelService) handleNewAvatar(event bus.Event) {
	payload := event.Payload()
	userID, _ := payload["user_id"].(string)
	newAvatar, _ := payload["new_avatar"].(string)

	if userID == "" || newAvatar == "" {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	err := s.channelRepository.UpdateUserAvatar(ctx, userID, newAvatar)
	if err != nil {
		log.Printf("Failed to update avatar: %v", err)
	}
}

func (s *channelService) CreateChannel(req *dto.CreateChannelRequest, requesterID string) (*model.Channel, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Require at least 2 members
	if len(req.Members) < 2 {
		return nil, apperror.ErrBadRequest
	}

	isMember := false
	for _, m := range req.Members {
		if m.ID.Hex() == requesterID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, apperror.ErrBadRequest
	}

	return s.channelRepository.Create(ctx, req, requesterID)
}

func (s *channelService) GetChannelByID(channelID string) (*model.Channel, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.channelRepository.GetByID(ctx, channelID)
}

func (s *channelService) GetChannelsByUserID(userID string, requesterID string, page int, pageSize int) (*dto.PaginatedChannelsResponse, error) {
	if userID != requesterID {
		return nil, apperror.ErrForbidden
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channels, total, err := s.channelRepository.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	channelResponses := dto.FromChannels(channels)

	var response = dto.PaginatedChannelsResponse{
		Channels: channelResponses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}

	return &response, nil
}

func (s *channelService) GetChannelByBothUserID(user1ID string, user2ID string, requesterID string) (*model.Channel, error) {
	if requesterID != user1ID && requesterID != user2ID {
		return nil, apperror.ErrForbidden
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.channelRepository.GetByBothUserID(ctx, user1ID, user2ID)
}

func (s *channelService) UpdateChannel(req *dto.UpdateChannelRequest, requesterID string) (*model.Channel, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channel, err := s.channelRepository.GetByID(ctx, req.ChannelID)
	if err != nil {
		return nil, err
	}

	isMember := false
	for _, member := range channel.Members {
		if member.ID.Hex() == requesterID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, err
	}

	updated := false
	for i := range channel.Settings {
		if channel.Settings[i].UserID == requesterObjectID {
			if req.Notification != nil {
				channel.Settings[i].Notification = *req.Notification
			}
			if req.TypingIndicator != nil {
				channel.Settings[i].TypingIndicator = *req.TypingIndicator
			}
			updated = true
			break
		}
	}

	if !updated {
		return nil, apperror.ErrNoFieldsToUpdate
	}

	channel.UpdatedAt = time.Now()

	return s.channelRepository.Update(ctx, channel)
}

func (s *channelService) DeleteChannel(channelID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.channelRepository.Delete(ctx, channelID, userID)
}
