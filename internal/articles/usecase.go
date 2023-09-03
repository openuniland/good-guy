package articles

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/models"
)

type UseCase interface {
	FindOne(ctx context.Context) (*models.Article, error)
	UpdatedWithNewArticle(ctx context.Context) (*types.UpdatedWithNewArticleResponse, error)
	AddNewSubscriber(ctx context.Context, id string) error
	RemoveSubscriber(ctx context.Context, id string) error
}
