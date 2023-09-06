package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/internal/middlewares"
)

func MapFacebookRoutes(facebookGroup *gin.RouterGroup, h facebook.Handlers, mw *middlewares.MiddlewareManager) {
	facebookGroup.POST("/messages/text/:id", mw.AdminMiddleware(), h.SendMessage())
	facebookGroup.POST("/messages/button/:id", mw.AdminMiddleware(), h.SendButtonMessage())
	facebookGroup.POST("/quick-replies/:id", mw.AdminMiddleware(), h.SendQuickReplies())
}
