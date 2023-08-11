package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/ctms"
)

func MapCtmsRoutes(ctmsGroup *gin.RouterGroup, h ctms.Handlers) {
	ctmsGroup.POST("/login", h.Login())
	ctmsGroup.POST("/logout", h.Logout())
	ctmsGroup.POST("/daily-schedules", h.GetDailySchedule())
	ctmsGroup.POST("/exam-schedules", h.GetExamSchedule())
	ctmsGroup.POST("/upcoming-exam-schedules", h.GetUpcomingExamSchedule())
}
