package cookies

import "github.com/gin-gonic/gin"

type Handlers interface {
	CreateNewCookie() gin.HandlerFunc
	FindOneCookie() gin.HandlerFunc
	UpdateCookie() gin.HandlerFunc
}
