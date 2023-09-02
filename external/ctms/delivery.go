package ctms

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	LoginCtms() gin.HandlerFunc
	LogoutCtms() gin.HandlerFunc
	GetDailySchedule() gin.HandlerFunc
	GetExamSchedule() gin.HandlerFunc
	GetUpcomingExamSchedule() gin.HandlerFunc
	SendChangedExamScheduleAndNewExamScheduleToUser() gin.HandlerFunc
}
