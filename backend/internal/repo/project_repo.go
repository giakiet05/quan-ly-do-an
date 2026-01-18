package repo

import (
	"context"
	"errors"
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
	CreateProjectRounds(ctx context.Context, classroomID string, rounds []model.ProjectRound) error
	GetProjectRoundByID(ctx context.Context, classroomID, roundID string) (*model.ProjectRound, error)
	GetProjectRoundsByClassroomID(ctx context.Context, classroomID string) ([]model.ProjectRound, error)
	ListProjectRounds(ctx context.Context, classroomID string, page, pageSize int) ([]model.ProjectRound, int64, error)
	DeleteProjectRound(ctx context.Context, classroomID, roundID string) error

	CreateProject(ctx context.Context, project *model.Project) error
	CreateProjects(ctx context.Context, projects []model.Project) error
	GetProjectByID(ctx context.Context, classroomID, projectID string) (*model.Project, error)
	GetProjectsByRoundID(ctx context.Context, classroomID, roundID string) ([]model.Project, error)
	ReplaceProject(ctx context.Context, project *model.Project) error
	DeleteProject(ctx context.Context, classroomID, projectID string) error

	CreateReportPeriod(ctx context.Context, classroomID, roundID string, period *model.ReportPeriod) error
	CreateReportPeriods(ctx context.Context, classroomID, roundID string, reportPeriods []model.ReportPeriod) error
	GetReportPeriodByID(ctx context.Context, classroomID, roundID, reportPeriodID string) (*model.ReportPeriod, error)
	GetReportPeriodsByRoundID(ctx context.Context, classroomID, roundID string) ([]model.ReportPeriod, error)
	ReplaceReportPeriod(ctx context.Context, reportPeriod *model.ReportPeriod) error
	DeleteReportPeriod(ctx context.Context, classroomID string, roundID string, reportPeriodID string) error
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

func (p *projectRepo) CreateProjectRounds(
	ctx context.Context,
	classroomID string,
	rounds []model.ProjectRound,
) error {

	if len(rounds) == 0 {
		return nil
	}

	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}

	now := time.Now()
	for i := range rounds {
		rounds[i].ID = primitive.NewObjectID()
		rounds[i].CreatedAt = now
	}

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{
		"$push": bson.M{
			"rounds": bson.M{
				"$each": rounds,
			},
		},
	}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return apperror.ErrClassroomNotFound
	}

	return nil
}

func (p *projectRepo) GetProjectRoundByID(ctx context.Context, classroomID, roundID string) (*model.ProjectRound, error) {
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

func (p *projectRepo) GetProjectRoundsByClassroomID(
	ctx context.Context,
	classroomID string,
) ([]model.ProjectRound, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	var result struct {
		Rounds []model.ProjectRound `bson:"rounds"`
	}

	err = p.classroomCollection.FindOne(
		ctx,
		bson.M{"_id": classroomOID},
		options.FindOne().SetProjection(bson.M{
			"rounds": 1,
		}),
	).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrClassroomNotFound
		}
		return nil, err
	}

	return result.Rounds, nil
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

func (p *projectRepo) DeleteProjectRound(
	ctx context.Context,
	classroomID, roundID string,
) error {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	roundOID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	filter := bson.M{
		"_id":        classroomOID,
		"rounds._id": roundOID,
	}

	update := bson.M{
		"$pull": bson.M{
			"rounds": bson.M{
				"_id": roundOID,
			},
		},
	}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return apperror.ErrRoundNotFound
	}

	return nil
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

func (p *projectRepo) GetProjectByID(
	ctx context.Context,
	classroomID string,
	projectID string,
) (*model.Project, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	projectOID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id": classroomOID,
		}}},
		{{Key: "$unwind", Value: "$rounds"}},
		{{Key: "$unwind", Value: "$rounds.projects"}},
		{{Key: "$match", Value: bson.M{
			"rounds.projects._id": projectOID,
		}}},
		{{Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$rounds.projects",
		}}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, apperror.ErrProjectNotFound
	}

	var project model.Project
	if err := cursor.Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (p *projectRepo) GetProjectsByRoundID(
	ctx context.Context,
	classroomID, roundID string,
) ([]model.Project, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	roundOID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	var result struct {
		Rounds []struct {
			Projects []model.Project `bson:"projects"`
		} `bson:"rounds"`
	}

	err = p.classroomCollection.FindOne(
		ctx,
		bson.M{"_id": classroomOID},
		options.FindOne().SetProjection(bson.M{
			"rounds": bson.M{
				"$filter": bson.M{
					"input": "$rounds",
					"as":    "r",
					"cond": bson.M{
						"$eq": []interface{}{"$$r._id", roundOID},
					},
				},
			},
		}),
	).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrClassroomNotFound
		}
		return nil, err
	}

	if len(result.Rounds) == 0 {
		return nil, apperror.ErrRoundNotFound
	}

	return result.Rounds[0].Projects, nil
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

func (p *projectRepo) DeleteProject(
	ctx context.Context,
	classroomID, projectID string,
) error {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	projectOID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	filter := bson.M{
		"_id":                 classroomOID,
		"rounds.projects._id": projectOID,
	}

	update := bson.M{
		"$pull": bson.M{
			"rounds.$[].projects": bson.M{
				"_id": projectOID,
			},
		},
	}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
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

func (p *projectRepo) CreateReportPeriods(
	ctx context.Context,
	classroomID, roundID string,
	reportPeriods []model.ReportPeriod,
) error {

	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}
	roundObjectID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":        classroomObjectID,
		"rounds._id": roundObjectID,
	}

	update := bson.M{
		"$push": bson.M{
			"rounds.$.report_periods": bson.M{
				"$each": reportPeriods,
			},
		},
	}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return apperror.ErrRoundNotFound
	}

	return nil
}

func (p *projectRepo) GetReportPeriodByID(
	ctx context.Context,
	classroomID, roundID, reportPeriodID string,
) (*model.ReportPeriod, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}
	roundOID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return nil, err
	}
	reportPeriodOID, err := primitive.ObjectIDFromHex(reportPeriodID)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id": classroomOID,
		}}},
		{{Key: "$unwind", Value: "$rounds"}},
		{{Key: "$match", Value: bson.M{
			"rounds._id": roundOID,
		}}},
		{{Key: "$unwind", Value: "$rounds.report_periods"}},
		{{Key: "$match", Value: bson.M{
			"rounds.report_periods._id": reportPeriodOID,
		}}},
		{{Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$rounds.report_periods",
		}}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, apperror.ErrNotFound
	}

	var period model.ReportPeriod
	if err := cursor.Decode(&period); err != nil {
		return nil, err
	}

	return &period, nil
}

func (p *projectRepo) GetReportPeriodsByRoundID(
	ctx context.Context,
	classroomID, roundID string,
) ([]model.ReportPeriod, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	roundOID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	var result struct {
		Rounds []struct {
			ReportPeriods []model.ReportPeriod `bson:"report_periods"`
		} `bson:"rounds"`
	}

	err = p.classroomCollection.FindOne(
		ctx,
		bson.M{"_id": classroomOID},
		options.FindOne().SetProjection(bson.M{
			"rounds": bson.M{
				"$filter": bson.M{
					"input": "$rounds",
					"as":    "r",
					"cond": bson.M{
						"$eq": []interface{}{"$$r._id", roundOID},
					},
				},
			},
		}),
	).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrClassroomNotFound
		}
		return nil, err
	}

	if len(result.Rounds) == 0 {
		return nil, apperror.ErrRoundNotFound
	}

	return result.Rounds[0].ReportPeriods, nil
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
