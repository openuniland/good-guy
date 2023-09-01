package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/common"
)

func MapCommonRoutes(commonGroup *gin.RouterGroup, h common.Handlers) {
	commonGroup.POST("/webhook", h.HandleFacebookWebhook())
	commonGroup.GET("/webhook", h.VerifyFacebookWebhook())
}
