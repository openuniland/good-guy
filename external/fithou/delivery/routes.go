package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/fithou"
)

func MapFithouRoutes(fithouGroup *gin.RouterGroup, h fithou.Handlers) {
	fithouGroup.GET("/articles/crawl", h.CrawlArticlesFromFirstPage())
}
