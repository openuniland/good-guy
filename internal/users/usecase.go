package users

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UseCase interface {
	CreateNewUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
}
