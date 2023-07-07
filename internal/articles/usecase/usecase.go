package usecase

import (
	"context"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/fithou"
	"github.com/openuniland/good-guy/external/types"
	articles "github.com/openuniland/good-guy/internal/articles"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type ArticleUS struct {
	cfg         *configs.Configs
	articleRepo articles.Repository
	fithouUS    fithou.UseCase
}

func NewArticleUseCase(cfg *configs.Configs, articleRepo articles.Repository, fithouUS fithou.UseCase) articles.UseCase {
	return &ArticleUS{cfg: cfg, articleRepo: articleRepo, fithouUS: fithouUS}
}

func (a *ArticleUS) FindOne(ctx context.Context) (*models.Article, error) {

	articles, err := a.articleRepo.FindOne(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("error while fetching articles")
		return nil, err
	}

	return articles, nil
}

func (a *ArticleUS) UpdatedWithNewArticle(ctx context.Context) ([]*types.ArticleCrawl, error) {

	article, err := a.FindOne(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error while fetching articles")
		return nil, err
	}

	articles, err := a.fithouUS.CrawlArticlesFromFirstPage(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error while fetching articles")
		return nil, err
	}

	if len(articles) == 0 {
		return nil, nil
	}

	aidOfNewArticle := articles[0].Aid
	aidOfOldArticle := article.Aid

	var index int

	if aidOfNewArticle != aidOfOldArticle+1 {
		for i := 0; i < len(articles); i++ {
			if articles[i].Aid == aidOfOldArticle+1 {
				index = i
				break
			}
		}
	}

	articles = articles[:index]

	a.articleRepo.FindOneAndUpdate(ctx, bson.D{
		{
			Key:   "_id",
			Value: article.Id,
		},
	}, articles[0])

	return articles, nil
}
