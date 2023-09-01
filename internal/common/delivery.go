package common

import "github.com/gin-gonic/gin"

type Handlers interface {
	VerifyFacebookWebhook() gin.HandlerFunc
	HandleFacebookWebhook() gin.HandlerFunc
}
