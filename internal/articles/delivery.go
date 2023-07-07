package articles

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	FindOne() gin.HandlerFunc
	UpdatedWithNewArticle() gin.HandlerFunc
}
