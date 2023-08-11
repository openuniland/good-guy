package examschedules

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UseCase interface {
	CreateNewExamSchedules(ctx context.Context, examSchedules *models.ExamSchedules) (*mongo.InsertOneResult, error)
	FindExamSchedules(ctx context.Context) ([]*models.ExamSchedules, error)
	FindExamSchedulesByUsername(ctx context.Context, filter interface{}) (*models.ExamSchedules, error)
	UpdateExamSchedulesByUsername(ctx context.Context, filter interface{}, examSchedules *models.ExamSchedules) (*models.ExamSchedules, error)
}
