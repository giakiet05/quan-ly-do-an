package bus

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// Event Topics
const (
	TopicBroadcast = "broadcast"

	TopicWSConnected    = "ws.connected"
	TopicWSDisconnected = "ws.disconnected"
	TopicWSPackageSend  = "ws.send"

	TopicNewMessage    = "message.new"
	TopicMessageError  = "message.error"
	TopicTypingMessage = "message.typing"
	TopicInChatMessage = "message.in_chat"

	TopicUserChangeAvatar = "user.avatar"

	// Notification events
	TopicNotificationCreated = "notification.created"
	TopicNotificationUpdated = "notification.updated"
	TopicNotificationDeleted = "notification.deleted"

	// Group
	TopicGroupInvitation = "group.invitation"

	// Report
	TopicReportOpened       = "report.opened"
	TopicReportNearDeadline = "report.near_deadline"
	TopicReportSubmitted    = "report.submitted"
	TopicReportGraded       = "report.graded"
	TopicReportExpired      = "report.expired"

	// Project
	TopicProjectRegistrationOpened   = "project.registration.opened"
	TopicProjectRegistrationDeadline = "project.registration.deadline"
	TopicProjectRegistrationExpired  = "project.registration.expired"

	// Class
	TopicClassUpdated         = "class.updated"
	TopicClassroomInvitation  = "classroom.invitation"
)

type BroadcastEventType string

const (
	// ---- Message-related ----
	BroadcastEventMessageCreated BroadcastEventType = "message_created"
	BroadcastEventTypingStart    BroadcastEventType = "typing_start"
	BroadcastEventTypingStop     BroadcastEventType = "typing_stop"
	BroadcastEventMessageRead    BroadcastEventType = "message_read"

	// ---- Notification-related ----
	BroadcastEventMessageNotification BroadcastEventType = "message_notification"
)

type BroadcastEvent struct {
	RecipientIDs []string           `json:"recipient_ids"`
	EventType    BroadcastEventType `json:"event_type"`
	TempID       string             `json:"temp_id"`
	Data         interface{}        `json:"data"`
}

func (e BroadcastEvent) Topic() string {
	return TopicBroadcast
}
func (e BroadcastEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"recipient_ids": e.RecipientIDs,
		"event_type":    e.EventType,
		"temp_id":       e.TempID,
		"data":          e.Data,
	}
}

type WSPackageSendEvent struct {
	Type dto.WebSocketMessageType `json:"type"`
	Data interface{}              `json:"data"`
}

func (e WSPackageSendEvent) Topic() string { return TopicWSPackageSend }
func (e WSPackageSendEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"type": e.Type,
		"data": e.Data,
	}
}

type WSConnectedEvent struct {
	UserID string `json:"user_id"`
}

func (e WSConnectedEvent) Topic() string { return TopicWSConnected }
func (e WSConnectedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"user_id": e.UserID,
	}
}

type WSDisconnectedEvent struct {
	UserID string `json:"user_id"`
}

func (e WSDisconnectedEvent) Topic() string { return TopicWSDisconnected }
func (e WSDisconnectedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"user_id": e.UserID,
	}
}

type UserChangeAvatarEventType struct {
	UserID    string
	NewAvatar string
}

func (e UserChangeAvatarEventType) Topic() string {
	return TopicUserChangeAvatar
}
func (e UserChangeAvatarEventType) Payload() map[string]interface{} {
	return map[string]interface{}{
		"user_id":    e.UserID,
		"new_avatar": e.NewAvatar,
	}
}

// --- Notification Events ---

type NotificationCreatedEvent struct {
	RecipientID  string
	Notification dto.NotificationResponse
}

func (e NotificationCreatedEvent) Topic() string {
	return TopicNotificationCreated
}
func (e NotificationCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"recipientId":  e.RecipientID,
		"notification": e.Notification,
	}
}

// --- Message Events ---

type NewMessageEvent struct {
	TempMessageID  string            `json:"temp_message_id"`
	ChannelID      string            `json:"channel_id"`
	SenderID       string            `json:"sender_id"`
	SenderUsername string            `json:"sender_username"`
	Type           model.MessageType `json:"type"`
	Content        string            `json:"content"`
}

func (e NewMessageEvent) Topic() string {
	return TopicNewMessage
}

func (e NewMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"temp_message_id": e.TempMessageID,
		"channel_id":      e.ChannelID,
		"sender_id":       e.SenderID,
		"sender_username": e.SenderUsername,
		"type":            e.Type,
		"content":         e.Content,
	}
}

type MessageErrorEvent struct {
	SenderID      string `json:"sender_id"`
	ChannelID     string `json:"channel_id"`
	TempMessageID string `json:"temp_message_id"`
	ErrorCode     string `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
}

func (e MessageErrorEvent) Topic() string {
	return TopicMessageError
}

func (e MessageErrorEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"sender_id":       e.SenderID,
		"channel_id":      e.ChannelID,
		"temp_message_id": e.TempMessageID,
		"error_code":      e.ErrorCode,
		"error_msg":       e.ErrorMsg,
	}
}

type TypingMessageEvent struct {
	ChannelID string `json:"channel_id"`
	SenderID  string `json:"sender_id"`
	IsTyping  bool   `json:"is_typing"`
}

func (e TypingMessageEvent) Topic() string {
	return TopicTypingMessage
}

func (e TypingMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"channel_id": e.ChannelID,
		"sender_id":  e.SenderID,
		"is_typing":  e.IsTyping,
	}
}

type InChatMessageEvent struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	IsInChat  bool   `json:"is_in_chat"`
}

func (e InChatMessageEvent) Topic() string {
	return TopicInChatMessage
}

func (e InChatMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"channel_id": e.ChannelID,
		"user_id":    e.UserID,
		"is_in_chat": e.IsInChat,
	}
}

type GroupInvitationEvent struct {
	GroupID     string     `json:"group_id"`
	InviterID   string     `json:"inviter_id"`
	InviteeID   string     `json:"invitee_id"`
	IsAccepted  bool       `json:"is_accepted"`
	SentAt      time.Time  `json:"sent_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
}

func (e GroupInvitationEvent) Topic() string {
	return TopicGroupInvitation
}
func (e GroupInvitationEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"group_id":     e.GroupID,
		"inviter_id":   e.InviterID,
		"invitee_id":   e.InviteeID,
		"is_accepted":  e.IsAccepted,
		"sent_at":      e.SentAt,
		"responded_at": e.RespondedAt,
	}
}

type ClassroomInvitationEvent struct {
	InvitationID  string `json:"invitation_id"`
	ClassroomID   string `json:"classroom_id"`
	ClassroomName string `json:"classroom_name"`
	InviterID     string `json:"inviter_id"`
	InviterName   string `json:"inviter_name"`
	InviteeID     string `json:"invitee_id"`
}

func (e ClassroomInvitationEvent) Topic() string {
	return TopicClassroomInvitation
}
func (e ClassroomInvitationEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"invitation_id":  e.InvitationID,
		"classroom_id":   e.ClassroomID,
		"classroom_name": e.ClassroomName,
		"inviter_id":     e.InviterID,
		"inviter_name":   e.InviterName,
		"invitee_id":     e.InviteeID,
//type ReportOpenedEvent struct {
//	ReportID string    `json:"report_id"`
//	OpenedAt time.Time `json:"opened_at"`
//}
//
//func (e ReportOpenedEvent) Topic() string {
//	return TopicReportOpened
//}
//func (e ReportOpenedEvent) Payload() map[string]interface{} {
//	return map[string]interface{}{
//		"report_id": e.ReportID,
//		"opened_at": e.OpenedAt,
//	}
//}

type TopicReportSubmittedEvent struct {
	ClassroomID string    `json:"classroom_id"`
	GroupID     string    `json:"group_id"`
	ReportID    string    `json:"report_id"`
	SubmittedAt time.Time `json:"submitted_at"`
}

func (e TopicReportSubmittedEvent) Topic() string {
	return TopicReportSubmitted
}

func (e TopicReportSubmittedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"classroom_id": e.ClassroomID,
		"group_id":     e.GroupID,
		"report_id":    e.ReportID,
		"submitted_at": e.SubmittedAt,
	}
}

type TopicReportGradedEvent struct {
	ClassroomID string    `json:"classroom_id"`
	GroupID     string    `json:"group_id"`
	ReportID    string    `json:"report_id"`
	GradedAt    time.Time `json:"graded_at"`
}

func (e TopicReportGradedEvent) Topic() string {
	return TopicReportGraded
}

func (e TopicReportGradedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"classroom_id": e.ClassroomID,
		"group_id":     e.GroupID,
		"report_id":    e.ReportID,
		"graded_at":    e.GradedAt,
	}
}
