package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/openuniland/good-guy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
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

func (u *userHandlers) GetUserBySubscribedId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subscribedId := ctx.Param("subscribed_id")

		user, err := u.userUC.GetUserBySubscribedId(ctx, subscribedId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user found",
			"data":    user,
		})
	}
}

func (u *userHandlers) FindOneAndUpdateUser() gin.HandlerFunc {
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

		filter := bson.M{"username": user.Username}

		updateUser := &models.User{
			Username:     user.Username,
			Password:     user.Password,
			SubscribedID: user.SubscribedID,
		}

		res, err := u.userUC.FindOneAndUpdateUser(ctx, filter, updateUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user updated",
			"data":    res,
		})
	}
}

func (u *userHandlers) FindOneAndDeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filter := ctx.Param("username")

		res, err := u.userUC.FindOneAndDeleteUser(ctx, filter)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user deleted",
			"data":    res,
		})
	}
}
