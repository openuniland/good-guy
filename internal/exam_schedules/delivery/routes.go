package delivery

import (
	"github.com/gin-gonic/gin"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
	"github.com/openuniland/good-guy/internal/middlewares"
)

func MapExamSchedulesRoutes(examSchedulesGroup *gin.RouterGroup, h examschedules.Handlers, mw *middlewares.MiddlewareManager) {
	examSchedulesGroup.POST("", mw.AdminMiddleware(), h.CreateNewExamSchedules())
	examSchedulesGroup.GET("", mw.AdminMiddleware(), h.FindExamSchedules())
	examSchedulesGroup.PATCH("", mw.AdminMiddleware(), h.UpdateExamSchedulesByUsername())
}
