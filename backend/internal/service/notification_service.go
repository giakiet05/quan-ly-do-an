package service

import (
	"context"
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
	classroomRepo    repo.ClassroomRepo
	groupRepo        repo.GroupRepo
	eventBus         *bus.EventBus
	redisClient      *redis.Client
}

func NewNotificationService(
	notificationRepo repo.NotificationRepo,
	userRepo repo.UserRepo,
	classroomRepo repo.ClassroomRepo,
	groupRepo repo.GroupRepo,
	bus *bus.EventBus,
	redis *redis.Client,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		classroomRepo:    classroomRepo,
		groupRepo:        groupRepo,
		eventBus:         bus,
		redisClient:      redis,
	}
}

func (s *notificationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.eventBus.Subscribe(bus.TopicBroadcast, eventChannel)
	s.eventBus.Subscribe(bus.TopicGroupInvitation, eventChannel)
	s.eventBus.Subscribe(bus.TopicGroupJoinRequest, eventChannel)
	s.eventBus.Subscribe(bus.TopicClassroomInvitation, eventChannel)
	s.eventBus.Subscribe(bus.TopicReportSubmitted, eventChannel)
	s.eventBus.Subscribe(bus.TopicReportGraded, eventChannel)

	log.Println("NotificationService started and subscribed to events.")

	go s.handleUnsentNotifications()

	go s.processEvents(eventChannel)
}

func (s *notificationService) handleUnsentNotifications() {}

func (s *notificationService) handleUnsentReportOpenedNotifications() {}

func (s *notificationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicBroadcast:
			s.handleBroadcast(event)
		case bus.TopicGroupInvitation:
			s.handleGroupInvitation(event)
		case bus.TopicGroupJoinRequest:
			s.handleGroupJoinRequest(event)
		case bus.TopicClassroomInvitation:
			s.handleClassroomInvitation(event)
		case bus.TopicReportSubmitted:
			s.handleReportSubmitted(event)
		case bus.TopicReportGraded:
			s.handleReportGraded(event)
		}
	}
}

func (s *notificationService) handleBroadcast(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()

	recipientIDs, _ := payload["recipient_ids"].([]string)
	if len(recipientIDs) == 0 {
		return
	}

	eventType, _ := payload["event_type"].(bus.BroadcastEventType)
	data := payload["data"]

	switch eventType {
	case bus.BroadcastEventMessageCreated:
		s.handleMessageCreated(ctx, recipientIDs, data)
	case bus.BroadcastEventReportOpened:
		s.handleReportOpened(ctx, recipientIDs, data)
	case bus.BroadcastReportNearDeadline:
		s.handleReportNearDeadline(ctx, recipientIDs, data)
	case bus.BroadcastEventProjectRegistrationOpened:
		s.handleProjectRegistrationOpened(ctx, recipientIDs, data)
	case bus.BroadcastEventProjectRegistrationDeadline:
		s.handleProjectRegistrationDeadline(ctx, recipientIDs, data)
	}
}

func (s *notificationService) handleMessageCreated(
	ctx context.Context,
	recipientIDs []string,
	data interface{},
) {
	var messageData dto.MessageResponse
	if err := util.DecodeJson(data, &messageData); err != nil {
		log.Printf("Failed to decode message: %v", err)
		return
	}

	senderOID, err := primitive.ObjectIDFromHex(messageData.SenderID)
	if err != nil {
		return
	}

	key := fmt.Sprintf(config.RedisActiveUsersKey, messageData.ChannelID)

	for _, rid := range recipientIDs {
		if rid == messageData.SenderID {
			continue
		}

		if s.isUserActiveInChannel(key, rid) {
			continue
		}

		recipientOID, err := primitive.ObjectIDFromHex(rid)
		if err != nil {
			continue
		}

		notification := &model.Notification{
			RecipientID: recipientOID,
			ActorID:     senderOID,
			Type:        model.NotificationTypeNewMessage,
			Message:     fmt.Sprintf("Tin nhắn mới từ %s", messageData.SenderUsername),
			Link:        fmt.Sprintf("/channels/%s", messageData.ChannelID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		s.createAndPublish(ctx, rid, notification)
	}
}

func (s *notificationService) handleReportOpened(
	ctx context.Context,
	recipientIDs []string,
	data interface{},
) {
	var reportData map[string]interface{}
	if err := util.DecodeJson(data, &reportData); err != nil {
		log.Printf("Failed to decode report opened data: %v", err)
		return
	}

	reportTitle, _ := reportData["report_period_title"].(string)
	projectRoundName, _ := reportData["project_round_name"].(string)
	classroomID, _ := reportData["classroom_id"].(string)
	reportPeriodID, _ := reportData["report_period_id"].(string)

	for _, rid := range recipientIDs {
		recipientOID, err := primitive.ObjectIDFromHex(rid)
		if err != nil {
			continue
		}

		notification := &model.Notification{
			RecipientID: recipientOID,
			Type:        model.NotificationTypeReportOpened,
			Message:     fmt.Sprintf("Đợt báo cáo '%s' của %s đã mở", reportTitle, projectRoundName),
			Link:        fmt.Sprintf("/classrooms/%s/reports/%s", classroomID, reportPeriodID),
			IsRead:      false,
			Metadata:    reportData,
			CreatedAt:   time.Now(),
		}

		s.createAndPublish(ctx, rid, notification)
	}
}

func (s *notificationService) handleReportNearDeadline(
	ctx context.Context,
	recipientIDs []string,
	data interface{},
) {
	var reportData map[string]interface{}
	if err := util.DecodeJson(data, &reportData); err != nil {
		log.Printf("Failed to decode report deadline data: %v", err)
		return
	}

	reportTitle, _ := reportData["report_period_title"].(string)
	daysRemaining, _ := reportData["days_remaining"].(int)
	classroomID, _ := reportData["classroom_id"].(string)
	reportPeriodID, _ := reportData["report_period_id"].(string)

	for _, rid := range recipientIDs {
		recipientOID, err := primitive.ObjectIDFromHex(rid)
		if err != nil {
			continue
		}

		notification := &model.Notification{
			RecipientID: recipientOID,
			Type:        model.NotificationTypeReportDeadline,
			Message: fmt.Sprintf("Báo cáo '%s' sắp đến hạn nộp (còn %d ngày)",
				reportTitle,
				daysRemaining,
			),
			Link:      fmt.Sprintf("/classrooms/%s/reports/%s", classroomID, reportPeriodID),
			IsRead:    false,
			Metadata:  reportData,
			CreatedAt: time.Now(),
		}

		s.createAndPublish(ctx, rid, notification)
	}
}

func (s *notificationService) handleProjectRegistrationOpened(
	ctx context.Context,
	recipientIDs []string,
	data interface{},
) {
	var projectData map[string]interface{}
	if err := util.DecodeJson(data, &projectData); err != nil {
		log.Printf("Failed to decode project registration opened data: %v", err)
		return
	}

	projectRoundName, _ := projectData["project_round_name"].(string)
	classroomID, _ := projectData["classroom_id"].(string)
	projectRoundID, _ := projectData["project_round_id"].(string)

	for _, rid := range recipientIDs {
		recipientOID, err := primitive.ObjectIDFromHex(rid)
		if err != nil {
			continue
		}

		notification := &model.Notification{
			RecipientID: recipientOID,
			Type:        model.NotificationTypeProjectRegistrationOpened,
			Message:     fmt.Sprintf("Đợt đăng ký đồ án '%s' đã mở", projectRoundName),
			Link:        fmt.Sprintf("/classrooms/%s/project-rounds/%s", classroomID, projectRoundID),
			IsRead:      false,
			Metadata:    projectData,
			CreatedAt:   time.Now(),
		}

		s.createAndPublish(ctx, rid, notification)
	}
}

func (s *notificationService) handleProjectRegistrationDeadline(
	ctx context.Context,
	recipientIDs []string,
	data interface{},
) {
	var projectData map[string]interface{}
	if err := util.DecodeJson(data, &projectData); err != nil {
		log.Printf("Failed to decode project registration deadline data: %v", err)
		return
	}

	projectRoundName, _ := projectData["project_round_name"].(string)
	daysRemaining, _ := projectData["days_remaining"].(int)
	classroomID, _ := projectData["classroom_id"].(string)
	projectRoundID, _ := projectData["project_round_id"].(string)

	for _, rid := range recipientIDs {
		recipientOID, err := primitive.ObjectIDFromHex(rid)
		if err != nil {
			continue
		}

		notification := &model.Notification{
			RecipientID: recipientOID,
			Type:        model.NotificationTypeProjectRegistrationDeadline,
			Message: fmt.Sprintf("Đợt đăng ký đồ án '%s' sắp kết thúc (còn %d ngày)",
				projectRoundName,
				daysRemaining,
			),
			Link:      fmt.Sprintf("/classrooms/%s/project-rounds/%s", classroomID, projectRoundID),
			IsRead:    false,
			Metadata:  projectData,
			CreatedAt: time.Now(),
		}

		s.createAndPublish(ctx, rid, notification)
	}
}

func (s *notificationService) isUserActiveInChannel(redisKey, userID string) bool {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isMember, err := s.redisClient.SIsMember(ctx, redisKey, userID).Result()
	if err != nil {
		log.Printf("Redis check failed: %v", err)
		return false
	}
	return isMember
}

func (s *notificationService) createAndPublish(
	ctx context.Context,
	recipientID string,
	notification *model.Notification,
) {
	created, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		return
	}

	s.eventBus.Publish(bus.NotificationCreatedEvent{
		RecipientID:  recipientID,
		Notification: dto.FromNotification(created),
	})
}

func (s *notificationService) handleGroupInvitation(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	inviteeID, _ := payload["invitee_id"].(string)
	inviterID, _ := payload["inviter_id"].(string)
	groupID, _ := payload["group_id"].(string)
	isAccepted, _ := payload["is_accepted"].(bool)
	respondedAt, _ := payload["responded_at"].(*time.Time)

	// If respondedAt is nil, it's a new invitation (sent)
	if respondedAt == nil {
		// Notify invitee about the invitation
		inviteeObjectID, _ := primitive.ObjectIDFromHex(inviteeID)
		inviterObjectID, _ := primitive.ObjectIDFromHex(inviterID)

		notification := &model.Notification{
			RecipientID: inviteeObjectID,
			ActorID:     inviterObjectID,
			Type:        model.NotificationTypeGroupInvitation,
			Message:     "Bạn có thư mời tham gia nhóm",
			Link:        fmt.Sprintf("/groups/%s", groupID),
			IsRead:      false,
			Metadata:    payload,
			CreatedAt:   time.Now(),
		}
		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  inviteeID,
			Notification: dto.FromNotification(createdNotification),
		})
	} else {
		// Invitation was responded to - notify inviter (group leader)
		inviterObjectID, _ := primitive.ObjectIDFromHex(inviterID)
		inviteeObjectID, _ := primitive.ObjectIDFromHex(inviteeID)

		var message string
		if isAccepted {
			message = "Lời mời tham gia nhóm đã được chấp nhận"
		} else {
			message = "Lời mời tham gia nhóm đã bị từ chối"
		}

		notification := &model.Notification{
			RecipientID: inviterObjectID,
			ActorID:     inviteeObjectID,
			Type:        model.NotificationTypeSystem,
			Message:     message,
			Link:        fmt.Sprintf("/groups/%s", groupID),
			IsRead:      false,
			Metadata:    payload,
			CreatedAt:   time.Now(),
		}
		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  inviterID,
			Notification: dto.FromNotification(createdNotification),
		})
	}
}

func (s *notificationService) handleGroupJoinRequest(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	requesterID, _ := payload["requester_id"].(string)
	leaderID, _ := payload["leader_id"].(string)
	groupID, _ := payload["group_id"].(string)
	status, _ := payload["status"].(string)
	updatedAt, _ := payload["updated_at"].(*time.Time)

	// If updatedAt is nil, it's a new request (pending)
	if updatedAt == nil {
		// Notify leader about the join request
		leaderObjectID, _ := primitive.ObjectIDFromHex(leaderID)
		requesterObjectID, _ := primitive.ObjectIDFromHex(requesterID)

		notification := &model.Notification{
			RecipientID: leaderObjectID,
			ActorID:     requesterObjectID,
			Type:        model.NotificationTypeSystem,
			Message:     "Bạn có yêu cầu tham gia nhóm mới",
			Link:        fmt.Sprintf("/groups/%s", groupID),
			IsRead:      false,
			Metadata:    payload,
			CreatedAt:   time.Now(),
		}
		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  leaderID,
			Notification: dto.FromNotification(createdNotification),
		})
	} else {
		// Request was responded to - notify requester
		requesterObjectID, _ := primitive.ObjectIDFromHex(requesterID)
		leaderObjectID, _ := primitive.ObjectIDFromHex(leaderID)

		var message string
		if status == string(model.RequestAccepted) {
			message = "Yêu cầu tham gia nhóm của bạn đã được chấp nhận"
		} else {
			message = "Yêu cầu tham gia nhóm của bạn đã bị từ chối"
		}

		notification := &model.Notification{
			RecipientID: requesterObjectID,
			ActorID:     leaderObjectID,
			Type:        model.NotificationTypeSystem,
			Message:     message,
			Link:        fmt.Sprintf("/groups/%s", groupID),
			IsRead:      false,
			Metadata:    payload,
			CreatedAt:   time.Now(),
		}
		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  requesterID,
			Notification: dto.FromNotification(createdNotification),
		})
	}
}

func (s *notificationService) handleClassroomInvitation(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	inviteeID, _ := payload["invitee_id"].(string)
	inviterID, _ := payload["inviter_id"].(string)
	inviterName, _ := payload["inviter_name"].(string)
	classroomName, _ := payload["classroom_name"].(string)
	classroomID, _ := payload["classroom_id"].(string)
	invitationID, _ := payload["invitation_id"].(string)

	inviteeObjectID, _ := primitive.ObjectIDFromHex(inviteeID)
	inviterObjectID, _ := primitive.ObjectIDFromHex(inviterID)

	notification := &model.Notification{
		RecipientID: inviteeObjectID,
		ActorID:     inviterObjectID,
		Type:        model.NotificationTypeClassroomInvitation,
		Message:     fmt.Sprintf("%s moi ban lam tro giang lop %s", inviterName, classroomName),
		Link:        fmt.Sprintf("/classrooms/%s", classroomID),
		IsRead:      false,
		Metadata: map[string]interface{}{
			"invitation_id":  invitationID,
			"classroom_id":   classroomID,
			"classroom_name": classroomName,
		},
		CreatedAt: time.Now(),
	}
	createdNotification, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		log.Printf("ERROR: NotificationService: failed to create classroom invitation notification: %v", err)
		return
	}

	s.eventBus.Publish(bus.NotificationCreatedEvent{
		RecipientID:  inviteeID,
		Notification: dto.FromNotification(createdNotification),
	})
}

func (s *notificationService) handleReportSubmitted(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	classroomID, _ := payload["classroom_id"].(string)
	groupID, _ := payload["group_id"].(string)
	reportID, _ := payload["report_id"].(string)
	submitterID, _ := payload["submitter_id"].(string)

	submitterOID, err := primitive.ObjectIDFromHex(submitterID)
	if err != nil {
		log.Printf("ERROR: Invalid submitter ID in report submitted event: %v", err)
		return
	}

	notification := &model.Notification{
		RecipientID: submitterOID,
		Type:        model.NotificationTypeReportSubmitted,
		Message:     "Bạn đã nộp báo cáo thành công",
		Link:        fmt.Sprintf("/classrooms/%s/groups/%s/reports/%s", classroomID, groupID, reportID),
		IsRead:      false,
		Metadata:    payload,
		CreatedAt:   time.Now(),
	}

	s.createAndPublish(ctx, submitterID, notification)
}

func (s *notificationService) handleReportGraded(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	classroomID, _ := payload["classroom_id"].(string)
	groupID, _ := payload["group_id"].(string)
	reportID, _ := payload["report_id"].(string)

	// Get group to get all members
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		log.Printf("ERROR: Failed to get group for report graded notification: %v", err)
		return
	}

	// Notify all group members
	for _, member := range group.Members {
		notification := &model.Notification{
			RecipientID: member.ID,
			Type:        model.NotificationTypeReportGraded,
			Message:     "Báo cáo của nhóm bạn đã được chấm điểm",
			Link:        fmt.Sprintf("/classrooms/%s/groups/%s/reports/%s", classroomID, groupID, reportID),
			IsRead:      false,
			Metadata:    payload,
			CreatedAt:   time.Now(),
		}

		s.createAndPublish(ctx, member.ID.Hex(), notification)
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
