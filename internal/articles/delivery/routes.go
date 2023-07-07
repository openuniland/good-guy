package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/articles"
)

func MapArticleRoutes(articleGroup *gin.RouterGroup, h articles.Handlers) {
	articleGroup.GET("/one", h.FindOne())
	articleGroup.PATCH("", h.UpdatedWithNewArticle())
}
