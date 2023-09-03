package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/auth"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/utils"
)

type authHandlers struct {
	cfg    *configs.Configs
	authUC auth.UseCase
}

func NewCtmsHandlers(cfg *configs.Configs, authUC auth.UseCase) auth.Handlers {
	return &authHandlers{
		cfg:    cfg,
		authUC: authUC,
	}
}

func (a *authHandlers) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &models.LoginRequest{}

		if err := ctx.ShouldBindJSON(req); err != nil {
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

		err = a.authUC.Login(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}
