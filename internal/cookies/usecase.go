package cookies

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UseCase interface {
	CreateNewCookie(ctx context.Context, cookie *models.Cookie) (*mongo.InsertOneResult, error)
	FindOneCookie(ctx context.Context, username string) (*models.Cookie, error)
	UpdateCookie(ctx context.Context, filter bson.M, update bson.M) error
}
