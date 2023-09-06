package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/articles"
	"github.com/openuniland/good-guy/internal/middlewares"
)

func MapArticleRoutes(articleGroup *gin.RouterGroup, h articles.Handlers, mw *middlewares.MiddlewareManager) {
	articleGroup.GET("/one", mw.AdminMiddleware(), h.FindOne())
	articleGroup.PATCH("", mw.AdminMiddleware(), h.UpdatedWithNewArticle())
}
