package service

import (
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
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

	s.eventBus.Subscribe(bus.TopicBroadcast, eventChannel)
	s.eventBus.Subscribe(bus.TopicGroupInvitation, eventChannel)

	log.Println("NotificationService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

func (s *notificationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicBroadcast:
			s.handleBroadcast(event)
		case bus.TopicGroupInvitation:
			s.handleGroupInvitation(event)
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

	if len(recipientIDs) == 0 {
		return
	}

	switch eventType {
	case bus.BroadcastEventMessageCreated:
		var messageData dto.MessageResponse
		if err := util.DecodeJson(data, &messageData); err != nil {
			log.Printf("Failed to decode message: %v", err)
			return
		}

		senderObjectID, err := primitive.ObjectIDFromHex(messageData.SenderID)
		if err != nil {
			return
		}

		// Redis key for active users
		key := fmt.Sprintf(config.RedisActiveUsersKey, messageData.ChannelID)
		for _, rid := range recipientIDs {
			if rid == messageData.SenderID {
				continue
			}

			// Check if user is currently active in this chat
			rctx, rcancel := util.NewDefaultDBContext()
			isInChat, err := s.redisClient.SIsMember(rctx, key, rid).Result()
			defer rcancel()
			if err != nil {
				log.Printf("Failed to check active users: %v", err)
				continue
			}
			if isInChat {
				continue
			}

			recipientObjectID, err := primitive.ObjectIDFromHex(rid)
			if err != nil {
				continue
			}

			notification := &model.Notification{
				RecipientID: recipientObjectID,
				ActorID:     senderObjectID,
				Type:        model.NotificationTypeNewMessage,
				Message:     fmt.Sprintf("Tin nhắn mới từ %s", messageData.SenderUsername),
				Link:        fmt.Sprintf("/channels/%s", messageData.ChannelID),
				IsRead:      false,
				CreatedAt:   time.Now(),
			}
			createdNotification, err := s.notificationRepo.Create(ctx, notification)
			if err != nil {
				log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
				continue
			}

			// Publish event for each recipient
			s.eventBus.Publish(bus.NotificationCreatedEvent{
				RecipientID:  rid,
				Notification: dto.FromNotification(createdNotification),
			})
		}
	}
}

func (s *notificationService) handleGroupInvitation(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	inviteeID, _ := payload["invitee_ids"].(string)
	inviterID, _ := payload["inviter_id"].(string)

	inviteeObjectID, _ := primitive.ObjectIDFromHex(inviteeID)
	inviterObjectID, _ := primitive.ObjectIDFromHex(inviterID)

	notification := &model.Notification{
		RecipientID: inviteeObjectID,
		ActorID:     inviterObjectID,
		Type:        model.NotificationTypeGroupInvitation,
		Message:     fmt.Sprintf("Bạn có thư mời tham gia nhóm"),
		Link:        "",
		IsRead:      false,
		Metadata:    payload,
		CreatedAt:   time.Now(),
	}
	createdNotification, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
	}

	s.eventBus.Publish(bus.NotificationCreatedEvent{
		RecipientID:  inviteeID,
		Notification: dto.FromNotification(createdNotification),
	})
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
