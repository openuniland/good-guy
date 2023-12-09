package hou

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	LoginHou() gin.HandlerFunc
	LogoutHou() gin.HandlerFunc
}
