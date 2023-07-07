package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/fithou"
)

type fithouHandlers struct {
	cfg      *configs.Configs
	fithouUC fithou.UseCase
}

func NewFithouHandlers(cfg *configs.Configs, fithouUC fithou.UseCase) fithou.Handlers {
	return &fithouHandlers{
		cfg:      cfg,
		fithouUC: fithouUC,
	}
}

func (c *fithouHandlers) CrawlArticlesFromFirstPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articles, err := c.fithouUC.CrawlArticlesFromFirstPage(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": articles})

	}
}
