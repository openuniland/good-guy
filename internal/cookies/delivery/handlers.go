package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/dtos"
	"github.com/openuniland/good-guy/internal/cookies"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cookieHandlers struct {
	cfg      *configs.Configs
	cookieUC cookies.UseCase
}

func NewCookieHandlers(cfg *configs.Configs, cookieUC cookies.UseCase) cookies.Handlers {
	return &cookieHandlers{
		cfg:      cfg,
		cookieUC: cookieUC,
	}
}

func (c *cookieHandlers) CreateNewCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie := &models.Cookie{}

		if err := ctx.ShouldBindJSON(&cookie); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := utils.ValidateStruct(ctx, cookie)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		_, err = c.cookieUC.CreateNewCookie(ctx, cookie)
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

func (c *cookieHandlers) FindOneCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		res, err := c.cookieUC.FindOneCookie(ctx, username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    res,
		})
	}
}

func (c *cookieHandlers) UpdateCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		cookie := &dtos.UpdateCookieRequest{}

		if err := ctx.ShouldBindJSON(&cookie); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := utils.ValidateStruct(ctx, cookie)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		filter := bson.M{"username": username}
		pushUpdate := bson.M{
			"$push": bson.M{"cookies": cookie.Cookie},
			"$set":  bson.M{"updated_at": primitive.NewDateTimeFromTime(time.Now())},
		}

		err = c.cookieUC.UpdateCookie(ctx, filter, pushUpdate)

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
