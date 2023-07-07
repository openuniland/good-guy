package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/articles"
)

type articleHandlers struct {
	cfg       *configs.Configs
	articleUC articles.UseCase
}

func NewArticleHandlers(cfg *configs.Configs, articleUC articles.UseCase) articles.Handlers {
	return &articleHandlers{
		cfg:       cfg,
		articleUC: articleUC,
	}
}

func (c *articleHandlers) FindOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		article, err := c.articleUC.FindOne(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, article)
	}
}

func (c *articleHandlers) UpdatedWithNewArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articles, err := c.articleUC.UpdatedWithNewArticle(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, articles)
	}
}
