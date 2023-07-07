package fithou

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	CrawlArticlesFromFirstPage(ctx context.Context) ([]*types.ArticleCrawl, error)
}
