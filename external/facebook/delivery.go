package facebook

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	SendMessage() gin.HandlerFunc
	SendButtonMessage() gin.HandlerFunc
	VerifyFacebookWebhook() gin.HandlerFunc
	HandleFacebookWebhook() gin.HandlerFunc
	SendQuickReplies() gin.HandlerFunc
}
