package repo

import (
	"context"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectRepo interface {
	CreateProjectRound(ctx context.Context, classroomID string, round *model.ProjectRound) error
	GetProjectRound(ctx context.Context, classroomID, roundID string) (*model.ProjectRound, error)
	ListProjectRounds(ctx context.Context, classroomID string, page, pageSize int) ([]model.ProjectRound, int64, error)

	CreateProject(ctx context.Context, project *model.Project) error
	CreateProjects(ctx context.Context, projects []model.Project) error
	ReplaceProject(ctx context.Context, project *model.Project) error
	DeleteProject(ctx context.Context, projectID string) error

	CreateReportPeriod(ctx context.Context, classroomID, roundID string, period *model.ReportPeriod) error
	ReplaceReportPeriod(ctx context.Context, reportPeriod *model.ReportPeriod) error
	DeleteReportPeriod(ctx context.Context, classroomID, roundID, reportPeriodID string) error
}

type projectRepo struct {
	classroomCollection *mongo.Collection
}

func NewProjectRepo(db *mongo.Database) ProjectRepo {
	return &projectRepo{classroomCollection: db.Collection(config.ClassroomColName)}
}

func (p *projectRepo) CreateProjectRound(ctx context.Context, classroomID string, round *model.ProjectRound) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}

	round.ID = primitive.NewObjectID()
	round.CreatedAt = time.Now()

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$push": bson.M{"rounds": round}}

	_, err = p.classroomCollection.UpdateOne(ctx, filter, update)
	return err
}

func (p *projectRepo) GetProjectRound(ctx context.Context, classroomID, roundID string) (*model.ProjectRound, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}
	roundObjectID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return nil, err
	}

	// Aggregation pipeline để lấy round cụ thể
	pipeline := []bson.M{
		{"$match": bson.M{"_id": classroomObjectID}},
		{"$unwind": "$rounds"},
		{"$match": bson.M{"rounds._id": roundObjectID}},
		{"$replaceRoot": bson.M{"newRoot": "$rounds"}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var round model.ProjectRound
	if !cursor.Next(ctx) {
		return nil, apperror.ErrRoundNotFound
	}

	err = cursor.Decode(&round)
	if err != nil {
		return nil, err
	}
	return &round, nil
}
func (p *projectRepo) ListProjectRounds(ctx context.Context, classroomID string, page, pageSize int) ([]model.ProjectRound, int64, error) {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize

	pipeline := []bson.M{
		{"$match": bson.M{"_id": classroomObjectID}},
		{"$unwind": "$rounds"},
		{"$sort": bson.D{{Key: "rounds.created_at", Value: -1}}},
		{"$skip": skip},
		{"$limit": pageSize},
		{"$replaceRoot": bson.M{"newRoot": "$rounds"}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var rounds []model.ProjectRound
	if err := cursor.All(ctx, &rounds); err != nil {
		return nil, 0, err
	}

	// COUNT TOTAL (pipeline riêng KHÔNG có skip/limit)
	countPipeline := []bson.M{
		{"$match": bson.M{"_id": classroomObjectID}},
		{"$unwind": "$rounds"},
		{"$count": "total"},
	}

	countCursor, err := p.classroomCollection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0]["total"].(int64)
	}

	return rounds, total, nil
}

func (p *projectRepo) CreateProject(ctx context.Context, project *model.Project) error {
	filter := bson.M{"_id": project.ClassroomID, "rounds._id": project.ProjectRoundID}
	update := bson.M{"$push": bson.M{"rounds.$.projects": project}}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrRoundNotFound
	}
	return nil
}

func (p *projectRepo) CreateProjects(ctx context.Context, projects []model.Project) error {
	for i := range projects {
		filter := bson.M{"_id": projects[i].ClassroomID, "rounds._id": projects[i].ProjectRoundID}
		update := bson.M{"$push": bson.M{"rounds.$.projects": projects[i]}}

		res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return apperror.ErrRoundNotFound
		}
	}
	return nil
}

func (p *projectRepo) ReplaceProject(ctx context.Context, project *model.Project) error {
	// find any classroom that contains this project
	filter := bson.M{"rounds.projects._id": project.ID}
	update := bson.M{"$set": bson.M{"rounds.$[r].projects.$[p]": project}}
	arrayFilters := options.ArrayFilters{Filters: []interface{}{
		bson.M{"r._id": project.ProjectRoundID},
		bson.M{"p._id": project.ID},
	}}
	opts := options.Update().SetArrayFilters(arrayFilters)

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrProjectNotFound
	}
	return nil
}

func (p *projectRepo) DeleteProject(ctx context.Context, projectID string) error {
	projectObjectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}

	filter := bson.M{"rounds.projects._id": projectObjectID}
	// remove from all rounds' projects arrays where present
	update := bson.M{"$pull": bson.M{"rounds.$[].projects": bson.M{"_id": projectObjectID}}}

	res, err := p.classroomCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrProjectNotFound
	}
	return nil
}

func (p *projectRepo) CreateReportPeriod(ctx context.Context, classroomID, roundID string, period *model.ReportPeriod) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}
	roundObjectID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": classroomObjectID, "rounds._id": roundObjectID}
	update := bson.M{"$push": bson.M{"rounds.$.report_periods": period}}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrRoundNotFound
	}
	return nil
}

func (p *projectRepo) ReplaceReportPeriod(ctx context.Context, reportPeriod *model.ReportPeriod) error {
	// find classroom/round that contains this report period and replace it
	filter := bson.M{"rounds.report_periods._id": reportPeriod.ID}
	update := bson.M{"$set": bson.M{"rounds.$[r].report_periods.$[p]": reportPeriod}}
	arrayFilters := options.ArrayFilters{Filters: []interface{}{
		bson.M{"p._id": reportPeriod.ID},
	}}
	opts := options.Update().SetArrayFilters(arrayFilters)

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (p *projectRepo) DeleteReportPeriod(ctx context.Context, classroomID, roundID, reportPeriodID string) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}
	roundObjectID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return err
	}
	reportPeriodObjectID, err := primitive.ObjectIDFromHex(reportPeriodID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": classroomObjectID, "rounds._id": roundObjectID}
	update := bson.M{"$pull": bson.M{"rounds.$.report_periods": bson.M{"_id": reportPeriodObjectID}}}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
