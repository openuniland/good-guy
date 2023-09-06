package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/fithou"
	"github.com/openuniland/good-guy/internal/middlewares"
)

func MapFithouRoutes(fithouGroup *gin.RouterGroup, h fithou.Handlers, mw *middlewares.MiddlewareManager) {
	fithouGroup.GET("/articles/crawl", mw.AdminMiddleware(), h.CrawlArticlesFromFirstPage())
}
