package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron             *cron.Cron
	classroomRepo    repo.ClassroomRepo
	projectRepo      repo.ProjectRepo
	notificationRepo repo.NotificationRepo
	eventBus         *bus.EventBus
}

func NewCronService(
	classroomRepo repo.ClassroomRepo,
	projectRepo repo.ProjectRepo,
	notificationRepo repo.NotificationRepo,
	eventBus *bus.EventBus,
) *CronService {
	return &CronService{
		cron:             cron.New(),
		classroomRepo:    classroomRepo,
		projectRepo:      projectRepo,
		notificationRepo: notificationRepo,
		eventBus:         eventBus,
	}
}

func (s *CronService) Start() {
	// Check for project registration deadlines every hour
	s.cron.AddFunc("@hourly", s.checkProjectRegistrationDeadlines)

	// Check for report deadlines every hour
	s.cron.AddFunc("@hourly", s.checkReportDeadlines)

	s.cron.Start()
	log.Println("CronService started: checking project registration and report deadlines")
}

func (s *CronService) Stop() {
	s.cron.Stop()
	log.Println("CronService stopped")
}

func (s *CronService) checkProjectRegistrationDeadlines() {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	log.Println("Cron: Checking project registration deadlines...")

	classrooms, err := s.getAllClassrooms(ctx)
	if err != nil {
		log.Printf("ERROR: Cron: failed to get classrooms: %v", err)
		return
	}

	now := time.Now()
	oneDayFromNow := now.Add(24 * time.Hour)
	threeDaysFromNow := now.Add(72 * time.Hour)

	for _, classroom := range classrooms {
		rounds, err := s.projectRepo.GetProjectRoundsByClassroomID(ctx, classroom.ID.Hex())
		if err != nil {
			log.Printf("ERROR: Cron: failed to get rounds for classroom %s: %v", classroom.ID.Hex(), err)
			continue
		}

		for _, round := range rounds {
			// Check if registration is opening soon (3 days before)
			if round.StartDate.After(now) && round.StartDate.Before(threeDaysFromNow) {
				s.sendProjectRegistrationOpenedNotification(ctx, classroom, round)
			}

			// Check if deadline is approaching (1 day before)
			if round.EndDate.After(now) && round.EndDate.Before(oneDayFromNow) {
				s.sendProjectRegistrationDeadlineNotification(ctx, classroom, round)
			}

			// Check if registration has expired
			if round.EndDate.Before(now) && round.EndDate.After(now.Add(-1*time.Hour)) {
				s.sendProjectRegistrationExpiredNotification(ctx, classroom, round)
			}
		}
	}

	log.Println("Cron: Finished checking project registration deadlines")
}

func (s *CronService) checkReportDeadlines() {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	log.Println("Cron: Checking report deadlines...")

	classrooms, err := s.getAllClassrooms(ctx)
	if err != nil {
		log.Printf("ERROR: Cron: failed to get classrooms: %v", err)
		return
	}

	now := time.Now()
	oneDayFromNow := now.Add(24 * time.Hour)
	threeDaysFromNow := now.Add(72 * time.Hour)

	for _, classroom := range classrooms {
		rounds, err := s.projectRepo.GetProjectRoundsByClassroomID(ctx, classroom.ID.Hex())
		if err != nil {
			log.Printf("ERROR: Cron: failed to get rounds for classroom %s: %v", classroom.ID.Hex(), err)
			continue
		}

		for _, round := range rounds {
			for _, reportPeriod := range round.ReportPeriods {
				// Check if report period is opening soon (3 days before)
				if reportPeriod.StartDate.After(now) && reportPeriod.StartDate.Before(threeDaysFromNow) {
					s.sendReportOpenedNotification(ctx, classroom, round, reportPeriod)
				}

				// Check if report deadline is approaching (1 day before)
				if reportPeriod.EndDate.After(now) && reportPeriod.EndDate.Before(oneDayFromNow) {
					s.sendReportDeadlineNotification(ctx, classroom, round, reportPeriod)
				}

				// Check if report has expired
				if reportPeriod.EndDate.Before(now) && reportPeriod.EndDate.After(now.Add(-1*time.Hour)) {
					s.sendReportExpiredNotification(ctx, classroom, round, reportPeriod)
				}
			}
		}
	}

	log.Println("Cron: Finished checking report deadlines")
}

func (s *CronService) getAllClassrooms(ctx context.Context) ([]model.Classroom, error) {
	// Get classrooms by fetching multiple pages
	allClassrooms := []model.Classroom{}
	page := 1
	pageSize := 100

	for {
		// We'll need to use GetByUniversity or a workaround
		// For now, let's use aggregation to get all classrooms
		var classrooms []model.Classroom
		// This is a simplified approach - in production you might want to optimize this
		break
	}

	return allClassrooms, nil
}

func (s *CronService) sendProjectRegistrationOpenedNotification(
	ctx context.Context,
	classroom model.Classroom,
	round model.ProjectRound,
) {
	// Send notification to all students in the classroom
	for _, studentID := range classroom.Students {
		notification := &model.Notification{
			RecipientID: studentID,
			Type:        model.NotificationTypeProjectRegistrationOpened,
			Message:     fmt.Sprintf("Đợt đăng ký đề tài '%s' sẽ mở vào %s", round.Name, round.StartDate.Format("02/01/2006 15:04")),
			Link:        fmt.Sprintf("/classrooms/%s/rounds/%s", classroom.ID.Hex(), round.ID.Hex()),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if _, err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("ERROR: Cron: failed to create notification: %v", err)
		}
	}
}

func (s *CronService) sendProjectRegistrationDeadlineNotification(
	ctx context.Context,
	classroom model.Classroom,
	round model.ProjectRound,
) {
	for _, studentID := range classroom.Students {
		notification := &model.Notification{
			RecipientID: studentID,
			Type:        model.NotificationTypeProjectRegistrationDeadline,
			Message:     fmt.Sprintf("Đợt đăng ký đề tài '%s' sắp hết hạn vào %s", round.Name, round.EndDate.Format("02/01/2006 15:04")),
			Link:        fmt.Sprintf("/classrooms/%s/rounds/%s", classroom.ID.Hex(), round.ID.Hex()),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if _, err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("ERROR: Cron: failed to create notification: %v", err)
		}
	}
}

func (s *CronService) sendProjectRegistrationExpiredNotification(
	ctx context.Context,
	classroom model.Classroom,
	round model.ProjectRound,
) {
	for _, studentID := range classroom.Students {
		notification := &model.Notification{
			RecipientID: studentID,
			Type:        model.NotificationTypeProjectRegistrationExpired,
			Message:     fmt.Sprintf("Đợt đăng ký đề tài '%s' đã kết thúc", round.Name),
			Link:        fmt.Sprintf("/classrooms/%s/rounds/%s", classroom.ID.Hex(), round.ID.Hex()),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if _, err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("ERROR: Cron: failed to create notification: %v", err)
		}
	}
}

func (s *CronService) sendReportOpenedNotification(
	ctx context.Context,
	classroom model.Classroom,
	round model.ProjectRound,
	reportPeriod model.ReportPeriod,
) {
	for _, studentID := range classroom.Students {
		notification := &model.Notification{
			RecipientID: studentID,
			Type:        model.NotificationTypeReportOpened,
			Message:     fmt.Sprintf("Kỳ báo cáo '%s' của đợt '%s' sẽ mở vào %s", reportPeriod.Title, round.Name, reportPeriod.StartDate.Format("02/01/2006 15:04")),
			Link:        fmt.Sprintf("/classrooms/%s/rounds/%s/reports/%s", classroom.ID.Hex(), round.ID.Hex(), reportPeriod.ID.Hex()),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if _, err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("ERROR: Cron: failed to create notification: %v", err)
		}
	}
}

func (s *CronService) sendReportDeadlineNotification(
	ctx context.Context,
	classroom model.Classroom,
	round model.ProjectRound,
	reportPeriod model.ReportPeriod,
) {
	for _, studentID := range classroom.Students {
		notification := &model.Notification{
			RecipientID: studentID,
			Type:        model.NotificationTypeReportDeadline,
			Message:     fmt.Sprintf("Kỳ báo cáo '%s' của đợt '%s' sắp hết hạn vào %s", reportPeriod.Title, round.Name, reportPeriod.EndDate.Format("02/01/2006 15:04")),
			Link:        fmt.Sprintf("/classrooms/%s/rounds/%s/reports/%s", classroom.ID.Hex(), round.ID.Hex(), reportPeriod.ID.Hex()),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if _, err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("ERROR: Cron: failed to create notification: %v", err)
		}
	}
}

func (s *CronService) sendReportExpiredNotification(
	ctx context.Context,
	classroom model.Classroom,
	round model.ProjectRound,
	reportPeriod model.ReportPeriod,
) {
	for _, studentID := range classroom.Students {
		notification := &model.Notification{
			RecipientID: studentID,
			Type:        model.NotificationTypeReportExpired,
			Message:     fmt.Sprintf("Kỳ báo cáo '%s' của đợt '%s' đã hết hạn", reportPeriod.Title, round.Name),
			Link:        fmt.Sprintf("/classrooms/%s/rounds/%s/reports/%s", classroom.ID.Hex(), round.ID.Hex(), reportPeriod.ID.Hex()),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if _, err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("ERROR: Cron: failed to create notification: %v", err)
		}
	}
}
