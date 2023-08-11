package examschedules

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, examSchedules *models.ExamSchedules) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}) ([]*models.ExamSchedules, error)
	FindOne(ctx context.Context, filter interface{}) (*models.ExamSchedules, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*models.ExamSchedules, error)
}
