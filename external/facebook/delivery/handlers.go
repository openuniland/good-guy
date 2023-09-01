package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/pkg/utils"
)

type facebookHandlers struct {
	cfg        *configs.Configs
	facebookUC facebook.UseCase
}

func NewFacebookHandlers(cfg *configs.Configs, facebookUC facebook.UseCase) facebook.Handlers {
	return &facebookHandlers{
		cfg:        cfg,
		facebookUC: facebookUC,
	}
}

func (f *facebookHandlers) SendMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		message := &types.SendMessageRequest{}

		if err := ctx.ShouldBindJSON(message); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err := utils.ValidateStruct(ctx, message)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		err = f.facebookUC.SendMessage(ctx, id, message)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Message sent!",
		})

	}
}

func (f *facebookHandlers) SendButtonMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		message := &types.SendButtonMessageRequest{}
		id := ctx.Param("id")

		if err := ctx.ShouldBindJSON(message); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err := utils.ValidateStruct(ctx, message)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		err = f.facebookUC.SendButtonMessage(ctx, id, message)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Message sent!",
		})

	}
}

func (f *facebookHandlers) SendQuickReplies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var req *types.SendQuickRepliesRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err := utils.ValidateStruct(ctx, req)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		err = f.facebookUC.SendQuickReplies(ctx, id, req.Text, &req.QuickReplies)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, "OK")
	}
}
