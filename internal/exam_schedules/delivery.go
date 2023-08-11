package examschedules

import "github.com/gin-gonic/gin"

type Handlers interface {
	CreateNewExamSchedules() gin.HandlerFunc
	FindExamSchedules() gin.HandlerFunc
	UpdateExamSchedulesByUsername() gin.HandlerFunc
}
