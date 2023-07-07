package articles

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, article *models.Article) (*mongo.InsertOneResult, error)
	Upsert(ctx context.Context, article *models.Article) (*mongo.UpdateResult, error)
	FindByAid(ctx context.Context, aid int) (*models.Article, error)
	Find(ctx context.Context, filter interface{}) ([]*models.Article, error)
	FindOne(ctx context.Context, filter interface{}) (*models.Article, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*models.Article, error)
}
