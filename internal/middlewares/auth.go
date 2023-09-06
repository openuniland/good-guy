package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader != mw.cfg.Server.AuthKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
