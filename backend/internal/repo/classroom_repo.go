package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomRepo interface {
	Create(ctx context.Context, classroom *model.Classroom) (*model.Classroom, error)
	GetByID(ctx context.Context, classroomID string) (*model.Classroom, error)
	Update(ctx context.Context, filter Filter, update UpdateDocument) error
	Replace(ctx context.Context, classroom *model.Classroom) error
	Delete(ctx context.Context, classroomID string) error

	IsLecturer(ctx context.Context, classroomID string, userID string) (bool, error)
}

type classroomRepo struct {
	classroomCollection *mongo.Collection
}

func NewClassroomRepo(db *mongo.Database) ClassroomRepo {
	return &classroomRepo{
		classroomCollection: db.Collection(config.ClassroomColName),
	}
}

func (c *classroomRepo) Create(ctx context.Context, classroom *model.Classroom) (*model.Classroom, error) {
	//TODO implement me
	panic("implement me")
}

func (c *classroomRepo) GetByID(ctx context.Context, classroomID string) (*model.Classroom, error) {
	//TODO implement me
	panic("implement me")
}

func (c *classroomRepo) Update(ctx context.Context, filter Filter, update UpdateDocument) error {
	//TODO implement me
	panic("implement me")
}

func (c *classroomRepo) Replace(ctx context.Context, classroom *model.Classroom) error {
	//TODO implement me
	panic("implement me")
}

func (c *classroomRepo) Delete(ctx context.Context, classroomID string) error {
	//TODO implement me
	panic("implement me")
}

func (c *classroomRepo) IsLecturer(ctx context.Context, classroomID string, userID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
