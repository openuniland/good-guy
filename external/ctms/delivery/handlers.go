package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/models"
	"github.com/openuniland/good-guy/pkg/utils"
)

type ctmsHandlers struct {
	cfg    *configs.Configs
	authUC ctms.UseCase
}

func NewCtmsHandlers(cfg *configs.Configs, authUC ctms.UseCase) ctms.Handlers {
	return &ctmsHandlers{
		cfg:    cfg,
		authUC: authUC,
	}
}

func (c *ctmsHandlers) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &models.LoginRequest{}

		if err := ctx.ShouldBindJSON(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err := utils.ValidateStruct(ctx, user)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		cookie, err := c.authUC.Login(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"cookie": cookie,
		})
	}
}
func (c *ctmsHandlers) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &models.LogoutRequest{}

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

		err = c.authUC.Logout(ctx, req.Cookie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}
