package fithou

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	CrawlArticlesFromFirstPage() gin.HandlerFunc
}
