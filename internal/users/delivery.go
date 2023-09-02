package users

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	CreateNewUser() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	GetUserBySubscribedId() gin.HandlerFunc
	FindOneAndUpdateUser() gin.HandlerFunc
	FindOneAndDeleteUser() gin.HandlerFunc
}
