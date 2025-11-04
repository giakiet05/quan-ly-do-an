package service

import (
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService interface {
	Start()
	GetNotifications(recipientID string, page, pageSize int) (*dto.PaginatedNotificationsResponse, error)
	MarkAllAsRead(recipientID string) (int64, error)
}

type notificationService struct {
	notificationRepo repo.NotificationRepo
	userRepo         repo.UserRepo
	eventBus         *bus.EventBus
	redisClient      *redis.Client
}

func NewNotificationService(notificationRepo repo.NotificationRepo, userRepo repo.UserRepo, bus *bus.EventBus, redis *redis.Client) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		eventBus:         bus,
		redisClient:      redis,
	}
}

func (s *notificationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	log.Println("NotificationService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

func (s *notificationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicBroadcast:
			s.handleBroadcast(event)
		}
	}
}

func (s *notificationService) handleBroadcast(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	recipientIDs, _ := payload["recipient_ids"].([]string)
	eventType, _ := payload["event_type"].(bus.BroadcastEventType)
	data := payload["data"]

	switch eventType {
	case bus.BroadcastEventMessageCreated:
		// Send notification to user not in chat
		var messageData dto.MessageResponse
		if err := util.DecodeJson(data, &messageData); err != nil {
			log.Printf("Failed to decode message: %v", err)
			return
		}

		key := fmt.Sprintf(config.RedisActiveUsersKey, messageData.ChannelID)
		isInChat, err := s.redisClient.SIsMember(ctx, key, recipientIDs[0]).Result()
		if err != nil {
			log.Printf("Failed to check active users: %v", err)
			return
		}

		if isInChat {
			// Recipient is in chat, skip notification
			return
		}

		recipientObjectID, err := primitive.ObjectIDFromHex(recipientIDs[0])
		if err != nil {
			return
		}

		actorObjectID, err := primitive.ObjectIDFromHex(messageData.SenderID)
		if err != nil {
			return
		}

		notification := &model.Notification{
			RecipientID: recipientObjectID,
			ActorID:     actorObjectID,
			Type:        model.NotificationTypeNewMessage,
			Message:     fmt.Sprintf("Tin nhắn mới từ %s", messageData.SenderUsername),
			Link:        fmt.Sprintf("/channels/%s", messageData.ChannelID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  recipientIDs[0],
			Notification: dto.FromNotification(createdNotification),
		})
	}
}

func (s *notificationService) GetNotifications(recipientID string, page, pageSize int) (*dto.PaginatedNotificationsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	notifications, total, err := s.notificationRepo.GetByRecipientID(ctx, recipientID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedNotificationsResponse{
		Notifications: dto.FromNotifications(notifications),
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}, nil
}

func (s *notificationService) MarkAllAsRead(recipientID string) (int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.notificationRepo.MarkAllAsRead(ctx, recipientID)
}
