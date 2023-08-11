package delivery

import (
	"github.com/gin-gonic/gin"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
)

func MapExamSchedulesRoutes(examSchedulesGroup *gin.RouterGroup, h examschedules.Handlers) {
	examSchedulesGroup.POST("", h.CreateNewExamSchedules())
	examSchedulesGroup.GET("", h.FindExamSchedules())
	examSchedulesGroup.PATCH("", h.UpdateExamSchedulesByUsername())
}
