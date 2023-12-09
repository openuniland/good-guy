package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/hou"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
)

type houHandlers struct {
	cfg   *configs.Configs
	houUC hou.UseCase
}

func NewHouHandlers(cfg *configs.Configs, houUC hou.UseCase) hou.Handlers {
	return &houHandlers{
		cfg:   cfg,
		houUC: houUC,
	}
}

func (h *houHandlers) LoginHou() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		req := &types.LoginHouRequest{}
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "should bind json",
			})
			return
		}

		err := utils.ValidateStruct(ctx, req)
		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		res, err := h.houUC.LoginHou(ctx, req)
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[HANDLERS]:[LoginHou]:[houUC.LoginHou]:[ERROR_INFO=%v, DATA=%v]", err, req)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error",
			})
			return
		}

		log.Info().Msgf("[INFO]:[HANDLERS]:[LoginHou]:[houUC.LoginHou]:[DATA=%v]", req)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    res,
		})
	}
}

func (h *houHandlers) LogoutHou() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h.houUC.LogoutHou(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Info().Msg("success")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}
