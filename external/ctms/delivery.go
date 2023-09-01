package ctms

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
	GetDailySchedule() gin.HandlerFunc
	GetExamSchedule() gin.HandlerFunc
	GetUpcomingExamSchedule() gin.HandlerFunc
	SendChangedExamScheduleAndNewExamScheduleToUser() gin.HandlerFunc
}
