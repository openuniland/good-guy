package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/internal/middlewares"
)

func MapCtmsRoutes(ctmsGroup *gin.RouterGroup, h ctms.Handlers, mw *middlewares.MiddlewareManager) {
	ctmsGroup.POST("/login", h.LoginCtms())
	ctmsGroup.POST("/logout", h.LogoutCtms())
	ctmsGroup.POST("/daily-schedules", mw.AdminMiddleware(), h.GetDailySchedule())
	ctmsGroup.POST("/exam-schedules", mw.AdminMiddleware(), h.GetExamSchedule())
	ctmsGroup.POST("/upcoming-exam-schedules", mw.AdminMiddleware(), h.GetUpcomingExamSchedule())
	ctmsGroup.POST("/exam-schedules/:id", mw.AdminMiddleware(), h.SendChangedExamScheduleAndNewExamScheduleToUser())
}
