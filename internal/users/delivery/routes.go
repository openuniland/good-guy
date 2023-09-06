package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/middlewares"
	"github.com/openuniland/good-guy/internal/users"
)

func MapUserRoutes(userGroup *gin.RouterGroup, h users.Handlers, mw *middlewares.MiddlewareManager) {
	userGroup.POST("", mw.AdminMiddleware(), h.CreateNewUser())
	userGroup.GET("", mw.AdminMiddleware(), h.GetUsers())
	userGroup.GET("/:subscribed_id", mw.AdminMiddleware(), h.GetUserBySubscribedId())
	userGroup.PUT("", mw.AdminMiddleware(), h.FindOneAndUpdateUser())
	userGroup.DELETE("/:username", mw.AdminMiddleware(), h.FindOneAndDeleteUser())
}
