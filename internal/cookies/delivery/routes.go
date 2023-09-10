package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/cookies"
	"github.com/openuniland/good-guy/internal/middlewares"
)

func MapCookieRoutes(cookieGroup *gin.RouterGroup, c cookies.Handlers, mw *middlewares.MiddlewareManager) {
	cookieGroup.POST("", mw.AdminMiddleware(), c.CreateNewCookie())
	cookieGroup.GET("/:username", mw.AdminMiddleware(), c.FindOneCookie())
	cookieGroup.PUT("/:username", mw.AdminMiddleware(), c.UpdateCookie())
}
