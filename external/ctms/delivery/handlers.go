package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/pkg/utils"
)

type ctmsHandlers struct {
	cfg    *configs.Configs
	ctmsUC ctms.UseCase
}

func NewCtmsHandlers(cfg *configs.Configs, ctmsUC ctms.UseCase) ctms.Handlers {
	return &ctmsHandlers{
		cfg:    cfg,
		ctmsUC: ctmsUC,
	}
}

func (c *ctmsHandlers) LoginCtms() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &types.LoginCtmsRequest{}

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

		cookie, err := c.ctmsUC.LoginCtms(ctx, user)
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
func (c *ctmsHandlers) LogoutCtms() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &types.LogoutCtmsRequest{}

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

		err = c.ctmsUC.LogoutCtms(ctx, req.Cookie)
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

func (c *ctmsHandlers) GetDailySchedule() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &types.GetDailyScheduleRequest{}

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

		data, err := c.ctmsUC.GetDailySchedule(ctx, req.Cookie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"data":    data,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    data,
		})
	}
}

func (c *ctmsHandlers) GetExamSchedule() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &types.GetExamScheduleRequest{}

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

		data, err := c.ctmsUC.GetExamSchedule(ctx, req.Cookie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"data":    data,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    data,
		})
	}
}

func (c *ctmsHandlers) GetUpcomingExamSchedule() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &types.LoginCtmsRequest{}

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

		data, err := c.ctmsUC.GetUpcomingExamSchedule(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"data":    data,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    data,
		})
	}
}

func (c *ctmsHandlers) SendChangedExamScheduleAndNewExamScheduleToUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &types.LoginCtmsRequest{}
		id := ctx.Param("id")

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

		err = c.ctmsUC.SendChangedExamScheduleAndNewExamScheduleToUser(ctx, req, id)
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
