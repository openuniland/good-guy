package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/facebook"
)

func MapFacebookRoutes(facebookGroup *gin.RouterGroup, h facebook.Handlers) {
	facebookGroup.POST("/messages/text/:id", h.SendMessage())
	facebookGroup.POST("/messages/button/:id", h.SendButtonMessage())
	facebookGroup.GET("/webhook", h.VerifyFacebookWebhook())
	facebookGroup.POST("/webhook", h.HandleFacebookWebhook())
}
