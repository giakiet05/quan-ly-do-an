package repo

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
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
	ReplaceProjectRound(ctx context.Context, round *model.ProjectRound) error
	DeleteProjectRound(ctx context.Context, classroomID, roundID string) error

	GetJustOpenReportPeriods(ctx context.Context, now time.Time) ([]dto.ReportPeriodWithProjectRound, error)
	GetNearDeadlineReportPeriods(ctx context.Context, now time.Time, daysBeforeDeadline int) ([]dto.ReportPeriodWithProjectRound, error)
	GetJustOpenProjectRounds(ctx context.Context, now time.Time) ([]dto.ProjectRoundWithClassroom, error)
	GetNearDeadlineProjectRounds(ctx context.Context, now time.Time, daysBeforeDeadline int) ([]dto.ProjectRoundWithClassroom, error)

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

	ReportPeriodExists(ctx context.Context, classroomID string, roundID string, periodID string) (bool, error)
}

type projectRepo struct {
	classroomCollection *mongo.Collection
	projectCollection   *mongo.Collection
}

func NewProjectRepo(db *mongo.Database) ProjectRepo {
	return &projectRepo{
		classroomCollection: db.Collection(config.ClassroomColName),
		projectCollection:   db.Collection(config.ProjectColName),
	}
}

func (p *projectRepo) CreateProjectRound(ctx context.Context, classroomID string, round *model.ProjectRound) error {
	classroomObjectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return err
	}

	round.CreatedAt = time.Now()

	filter := bson.M{"_id": classroomObjectID}
	update := bson.M{"$push": bson.M{"project_rounds": round}}

	result, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return apperror.ErrInternal
	}
	return nil
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

func (p *projectRepo) ReplaceProjectRound(ctx context.Context, round *model.ProjectRound) error {
	filter := bson.M{"project_rounds._id": round.ID}
	update := bson.M{"$set": bson.M{"project_rounds.$": round}}

	res, err := p.classroomCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return apperror.ErrRoundNotFound
	}
	return nil
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

func (p *projectRepo) GetJustOpenReportPeriods(ctx context.Context, now time.Time) ([]dto.ReportPeriodWithProjectRound, error) {
	loc := now.Location()

	startOfDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	)
	endOfDay := startOfDay.Add(24 * time.Hour)

	pipeline := mongo.Pipeline{
		// 1️⃣ explode rounds
		{{Key: "$unwind", Value: "$rounds"}},

		// 2️⃣ explode report_periods
		{{Key: "$unwind", Value: "$rounds.report_periods"}},

		// 3️⃣ filter by day + not deleted
		{{Key: "$match", Value: bson.M{
			"rounds.report_periods.start_date": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
			"rounds.is_deleted": false,
		}}},

		// 4️⃣ shape result
		{{Key: "$project", Value: bson.M{
			"_id":                0,
			"classroom_id":       "$_id",
			"project_round_id":   "$rounds._id",
			"project_round_name": "$rounds.name",
			"report_periods":     "$rounds.report_periods",
		}}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.ReportPeriodWithProjectRound
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *projectRepo) GetNearDeadlineReportPeriods(ctx context.Context, now time.Time, daysBeforeDeadline int) ([]dto.ReportPeriodWithProjectRound, error) {
	loc := now.Location()

	startOfToday := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	)
	endOfDeadlineDay := startOfToday.AddDate(0, 0, daysBeforeDeadline).Add(24 * time.Hour)

	pipeline := mongo.Pipeline{
		// 1️⃣ explode rounds
		{{Key: "$unwind", Value: "$rounds"}},

		// 2️⃣ explode report_periods
		{{Key: "$unwind", Value: "$rounds.report_periods"}},

		// 3️⃣ filter by deadline range + not deleted
		{{Key: "$match", Value: bson.M{
			"rounds.report_periods.end_date": bson.M{
				"$gte": startOfToday,
				"$lte": endOfDeadlineDay,
			},
			"rounds.is_deleted": false,
		}}},

		// 4️⃣ shape result
		{{Key: "$project", Value: bson.M{
			"_id":                0,
			"classroom_id":       "$_id",
			"project_round_id":   "$rounds._id",
			"project_round_name": "$rounds.name",
			"report_periods":     "$rounds.report_periods",
		}}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.ReportPeriodWithProjectRound
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *projectRepo) GetJustOpenProjectRounds(
	ctx context.Context,
	now time.Time,
) ([]dto.ProjectRoundWithClassroom, error) {

	loc := now.Location()

	startOfDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	)
	endOfDay := startOfDay.Add(24 * time.Hour)

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$rounds"}},

		{{Key: "$match", Value: bson.M{
			"rounds.start_date": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
			"rounds.is_deleted": false,
		}}},

		{{Key: "$project", Value: bson.M{
			"_id":            0,
			"classroom_id":   "$_id",
			"classroom_name": "$name",
			"round":          "$rounds",
		}}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.ProjectRoundWithClassroom
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *projectRepo) GetNearDeadlineProjectRounds(ctx context.Context, now time.Time, daysBeforeDeadline int) ([]dto.ProjectRoundWithClassroom, error) {
	loc := now.Location()

	startOfToday := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	)
	endOfDeadlineDay := startOfToday.AddDate(0, 0, daysBeforeDeadline).Add(24 * time.Hour)

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$rounds"}},

		{{Key: "$match", Value: bson.M{
			"rounds.end_date": bson.M{
				"$gte": startOfToday,
				"$lte": endOfDeadlineDay,
			},
			"rounds.is_deleted": false,
		}}},

		{{Key: "$project", Value: bson.M{
			"_id":            0,
			"classroom_id":   "$_id",
			"classroom_name": "$name",
			"round":          "$rounds",
		}}},
	}

	cursor, err := p.classroomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.ProjectRoundWithClassroom
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *projectRepo) CreateProject(ctx context.Context, project *model.Project) error {
	project.CreatedAt = time.Now()
	_, err := p.projectCollection.InsertOne(ctx, project)
	return err
}

func (p *projectRepo) CreateProjects(ctx context.Context, projects []model.Project) error {
	if len(projects) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(projects))
	for i := range projects {
		projects[i].CreatedAt = now
		docs[i] = projects[i]
	}

	_, err := p.projectCollection.InsertMany(ctx, docs)
	return err
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

	filter := bson.M{
		"_id":          projectOID,
		"classroom_id": classroomOID,
	}

	var project model.Project
	err = p.projectCollection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrProjectNotFound
		}
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

	filter := bson.M{
		"classroom_id":     classroomOID,
		"project_round_id": roundOID,
	}

	cursor, err := p.projectCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []model.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *projectRepo) ReplaceProject(ctx context.Context, project *model.Project) error {
	filter := bson.M{"_id": project.ID}

	res, err := p.projectCollection.ReplaceOne(ctx, filter, project)
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
		"_id":          projectOID,
		"classroom_id": classroomOID,
	}

	res, err := p.projectCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
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

func (p *projectRepo) ReportPeriodExists(
	ctx context.Context,
	classroomID string,
	roundID string,
	periodID string,
) (bool, error) {

	classroomOID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return false, err
	}

	roundOID, err := primitive.ObjectIDFromHex(roundID)
	if err != nil {
		return false, err
	}

	periodOID, err := primitive.ObjectIDFromHex(periodID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"_id": classroomOID,
		"rounds": bson.M{
			"$elemMatch": bson.M{
				"_id":                roundOID,
				"report_periods._id": periodOID,
			},
		},
	}

	count, err := p.classroomCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
