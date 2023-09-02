package users

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UseCase interface {
	CreateNewUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetUserBySubscribedId(ctx context.Context, subscribedId string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	FindOneAndUpdateUser(ctx context.Context, filter, update interface{}) (*models.User, error)
	FindOneAndDeleteUser(ctx context.Context, filter interface{}) (*mongo.SingleResult, error)
}
