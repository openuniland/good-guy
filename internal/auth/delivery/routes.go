package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/internal/auth"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers) {
	authGroup.POST("/login", h.Login())
}
