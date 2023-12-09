package http

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/external/hou"
)

func MapHouRoutes(houGroup *gin.RouterGroup, h hou.Handlers) {
	houGroup.POST("/login", h.LoginHou())
	houGroup.POST("/logout", h.LogoutHou())
}
