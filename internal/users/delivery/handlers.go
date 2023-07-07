package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/openuniland/good-guy/pkg/utils"
)

type userHandlers struct {
	cfg    *configs.Configs
	userUC users.UseCase
}

func NewArticleHandlers(cfg *configs.Configs, userUC users.UseCase) users.Handlers {
	return &userHandlers{
		cfg:    cfg,
		userUC: userUC,
	}
}

func (u *userHandlers) CreateNewUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &models.User{}

		if err := ctx.ShouldBindJSON(user); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := utils.ValidateStruct(ctx, user)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		res, err := u.userUC.CreateNewUser(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user created",
			"data":    res,
		})
	}
}

func (u *userHandlers) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := u.userUC.GetUsers(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "users found",
			"data":    users,
		})
	}
}
