package cookies

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, cookie *models.Cookie) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter bson.M) (*models.Cookie, error)
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error)
}
