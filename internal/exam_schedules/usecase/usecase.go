package usecase

import (
	"context"

	"github.com/openuniland/good-guy/configs"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExamSchedulesUS struct {
	cfg               *configs.Configs
	examschedulesRepo examschedules.Repository
}

func NewExamSchedulesUseCase(cfg *configs.Configs, examschedulesRepo examschedules.Repository) examschedules.UseCase {
	return &ExamSchedulesUS{cfg: cfg, examschedulesRepo: examschedulesRepo}
}

func (e *ExamSchedulesUS) CreateNewExamSchedules(ctx context.Context, examSchedules *models.ExamSchedules) (*mongo.InsertOneResult, error) {

	insertOneResult, err := e.examschedulesRepo.Create(ctx, examSchedules)
	if err != nil {
		return nil, err
	}

	return insertOneResult, nil
}

func (e *ExamSchedulesUS) FindExamSchedules(ctx context.Context) ([]*models.ExamSchedules, error) {

	examSchedules, err := e.examschedulesRepo.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return examSchedules, nil
}

func (e *ExamSchedulesUS) FindExamSchedulesByUsername(ctx context.Context, filter interface{}) (*models.ExamSchedules, error) {

	examSchedules, err := e.examschedulesRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return examSchedules, nil
}

func (e *ExamSchedulesUS) UpdateExamSchedulesByUsername(ctx context.Context, filter interface{}, examSchedules *models.ExamSchedules) (*models.ExamSchedules, error) {

	examSchedules, err := e.examschedulesRepo.FindOneAndUpdate(ctx, filter, examSchedules)

	if err != nil {
		return nil, err
	}

	return examSchedules, nil
}
