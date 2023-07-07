package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/users"
)

func MapUserRoutes(userGroup *gin.RouterGroup, h users.Handlers) {
	userGroup.POST("", h.CreateNewUser())
	userGroup.GET("", h.GetUsers())
}
