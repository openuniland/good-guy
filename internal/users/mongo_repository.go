package users

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	Find(ctx context.Context, filter interface{}) ([]*models.User, error)
	FindOneUserByCondition(ctx context.Context, filter interface{}) (*models.User, error)
	FindOneAndUpdate(ctx context.Context, filter, update bson.M) (*models.User, error)
	FindOneAndDelete(ctx context.Context, filter interface{}) *mongo.SingleResult
}
