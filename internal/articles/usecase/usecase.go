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

func (a *ArticleUS) UpdatedWithNewArticle(ctx context.Context) (*types.UpdatedWithNewArticleResponse, error) {

	oldArticle, err := a.FindOne(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error while fetching articles")
		return nil, err
	}

	articles, err := a.fithouUS.CrawlArticlesFromFirstPage(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error while fetching articles")
		return nil, err
	}

	if len(articles) <= 0 {
		return nil, nil
	}

	aidOfNewArticle := articles[0].Aid
	aidOfOldArticle := oldArticle.Aid

	var index int
	if articles[len(articles)-1].Aid > aidOfOldArticle {
		index = len(articles) - 1
	}

	if aidOfNewArticle != aidOfOldArticle {
		for i := 0; i < len(articles); i++ {
			if articles[i].Aid == aidOfOldArticle {
				index = i
				break
			}
		}
	}

	if index != 0 {
		articles = articles[:index]
	}
	if index == 0 {
		return &types.UpdatedWithNewArticleResponse{
			Data:  articles[:1],
			IsNew: false,
		}, nil
	}

	_, err = a.articleRepo.FindOneAndUpdate(ctx, bson.D{
		{
			Key:   "_id",
			Value: oldArticle.Id,
		},
	}, articles[0])

	if err != nil {
		log.Error().Err(err).Msg("error while fetching articles")
		return nil, err
	}

	return &types.UpdatedWithNewArticleResponse{
		Data:  articles,
		IsNew: true,
	}, nil
}
