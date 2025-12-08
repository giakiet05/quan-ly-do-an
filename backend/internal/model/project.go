package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegistrationRound struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`         // "Đợt 1", "Đợt 2"
	Status      RoundStatus        `bson:"status" json:"status"`     // open/closed
	Deadline    time.Time          `bson:"deadline" json:"deadline"` // Hạn chót đăng ký
	Description string             `bson:"description" json:"description"`
	Projects    []RoundProject     `bson:"projects" json:"projects"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type RoundStatus string

const (
	RoundOpen   RoundStatus = "open"
	RoundClosed RoundStatus = "closed"
)

// ==================== PROJECT TRONG ROUND ====================
type RoundProject struct {
	ID                     primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	Title                  string                  `bson:"title" json:"title"`
	Desc                   string                  `bson:"desc" json:"desc"`
	CurrentTeams           int                     `bson:"current_teams" json:"current_teams"`
	MaxTeams               int                     `bson:"max_teams" json:"max_teams"`
	MinMember              int                     `bson:"min_member" json:"min_member"`
	MaxMember              int                     `bson:"max_member" json:"max_member"`
	Teams                  []Team                  `bson:"teams" json:"teams"`
	GroupFormationRequests []GroupFormationRequest `bson:"group_formation_requests" json:"group_formation_requests"`
}

type Team struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Leader         UserInfo           `bson:"leader" json:"leader"`
	Members        []TeamMember       `bson:"members" json:"members"`
	JoinRequests   []JoinRequest      `bson:"join_requests" json:"join_requests"`
	CanJoinRequest bool               `bson:"can_join_request" json:"can_join_request"`
	MaxMember      int                `bson:"max_member" json:"max_member"`
	MinMember      int                `bson:"min_member" json:"min_member"`
	CurrentMember  int                `bson:"current_member" json:"current_member"`
	Status         TeamStatus         `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}
type GroupFormationRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LeaderID   primitive.ObjectID `bson:"leader_id" json:"leader_id"`
	LeaderInfo UserInfo           `bson:"leader_info" json:"leader_info"`
	Members    []FormationMember  `bson:"members" json:"members"`
	Status     FormationStatus    `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type FormationMember struct {
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserInfo UserInfo           `bson:"user_info" json:"user_info"`
	Status   FormationStatus    `bson:"status" json:"status"`
}
type FormationStatus string

const (
	FormationPending   FormationStatus = "pending"   // Chờ confirm
	FormationConfirmed FormationStatus = "confirmed" // Đủ min_member → Tạo team
	FormationRejected  FormationStatus = "rejected"  // Không đủ hoặc hủy
)

type TeamStatus string

const (
	TeamOpen   TeamStatus = "open"   // Đang tuyển
	TeamFull   TeamStatus = "full"   // Đủ người
	TeamClosed TeamStatus = "closed" // Đóng tuyển
)

type TeamMember struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserInfo UserInfo           `bson:"user_info" json:"user_info"`
	Role     MemberRole         `bson:"role" json:"role"`
	JoinedAt time.Time          `bson:"joined_at" json:"joined_at"`
}

type MemberRole string

const (
	LeaderTeam MemberRole = "leader"
	MemberTeam MemberRole = "member"
)

type JoinRequest struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserInfo    UserInfo           `bson:"user_info" json:"user_info"`
	Status      RequestStatus      `bson:"status" json:"status"`
	Message     string             `bson:"message" json:"message"` // Lý do muốn join
	RequestedAt time.Time          `bson:"requested_at" json:"requested_at"`
}

type RequestStatus string

const (
	RequestPending  RequestStatus = "pending"
	RequestAccepted RequestStatus = "accepted"
	RequestRejected RequestStatus = "rejected"
)
