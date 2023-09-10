package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/common"
)

type commonHandlers struct {
	cfg      *configs.Configs
	commonUC common.UseCase
}

func NewCommonHandlers(cfg *configs.Configs, commonUC common.UseCase) common.Handlers {
	return &commonHandlers{
		cfg:      cfg,
		commonUC: commonUC,
	}
}

func (c *commonHandlers) HandleFacebookWebhook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var message *types.FacebookWebhookRequest

		if err := ctx.ShouldBindJSON(&message); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err := c.commonUC.HandleFacebookWebhook(ctx, message)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.Writer.Write([]byte("EVENT_RECEIVED"))
	}
}

func (c *commonHandlers) VerifyFacebookWebhook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("hub.verify_token")
		challenge := ctx.Query("hub.challenge")

		res, err := c.commonUC.VerifyFacebookWebhook(ctx, token, challenge)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Forbidden",
			})
			return
		}

		ctx.Writer.Write([]byte(res))

	}
}
