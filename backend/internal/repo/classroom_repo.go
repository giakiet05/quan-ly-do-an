package repo

import (
	"context"

	"github.com/giakiet05/quan-ly-do-an/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomRepository interface {
	// Classroom
	CreateClassroom(ctx context.Context, classroom *model.Classroom) error
	GetClassroomByID(ctx context.Context, id primitive.ObjectID) (*model.Classroom, error)
	GetClassroomsByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*model.Classroom, error)
	AddStudentToClassroom(ctx context.Context, classroomID primitive.ObjectID, student model.UserInfo) error

	// Rounds
	CreateRound(ctx context.Context, classroomID primitive.ObjectID, round *model.RegistrationRound) error
	GetRound(ctx context.Context, classroomID, roundID primitive.ObjectID) (*model.RegistrationRound, error)
	UpdateRoundStatus(ctx context.Context, classroomID, roundID primitive.ObjectID, status model.RoundStatus) error

	// Projects
	CreateProject(ctx context.Context, classroomID, roundID primitive.ObjectID, project *model.RoundProject) error
	GetProject(ctx context.Context, classroomID, roundID, projectID primitive.ObjectID) (*model.RoundProject, error)
	UpdateProject(ctx context.Context, classroomID, roundID, projectID primitive.ObjectID, project *model.RoundProject) error

	// Group Formation
	CreateGroupFormationRequest(ctx context.Context, classroomID, roundID, projectID primitive.ObjectID, req *model.GroupFormationRequest) error
	ConfirmFormationMember(ctx context.Context, classroomID, roundID, projectID, formationID primitive.ObjectID, userID primitive.ObjectID, status model.FormationStatus) error
	CreateTeamFromFormation(ctx context.Context, classroomID, roundID, projectID, formationID primitive.ObjectID) (*model.Team, error)
	RejectFormationRequest(ctx context.Context, classroomID, roundID, projectID, formationID primitive.ObjectID) error

	// Teams
	CreateTeam(ctx context.Context, classroomID, roundID, projectID primitive.ObjectID, team *model.Team) error
	GetTeam(ctx context.Context, classroomID, roundID, projectID, teamID primitive.ObjectID) (*model.Team, error)
	UpdateTeam(ctx context.Context, classroomID, roundID, projectID, teamID primitive.ObjectID, team *model.Team) error
	AddJoinRequest(ctx context.Context, classroomID, roundID, projectID, teamID primitive.ObjectID, req *model.JoinRequest) error
	ApproveJoinRequest(ctx context.Context, classroomID, roundID, projectID, teamID, reqID primitive.ObjectID) error
	RejectJoinRequest(ctx context.Context, classroomID, roundID, projectID, teamID, reqID primitive.ObjectID) error
	MemberLeaveTeam(ctx context.Context, classroomID, roundID, projectID, teamID, userID primitive.ObjectID) error
	KickMember(ctx context.Context, classroomID, roundID, projectID, teamID, targetUserID primitive.ObjectID) error
	TransferLeadership(ctx context.Context, classroomID, roundID, projectID, teamID, newLeaderID primitive.ObjectID) error
	CheckAutoDissolve(ctx context.Context, classroomID, roundID, projectID, teamID primitive.ObjectID) error

	// Stats
	GetProjectStats(ctx context.Context, classroomID, roundID, projectID primitive.ObjectID) (*model.ProjectStats, error)
	GetRoundStats(ctx context.Context, classroomID, roundID primitive.ObjectID) (*model.RoundStats, error)
}
