package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type examSchedulesHandlers struct {
	cfg             *configs.Configs
	examschedulesUS examschedules.UseCase
}

func NewExamSchedulesHandlers(cfg *configs.Configs, examschedulesUS examschedules.UseCase) examschedules.Handlers {
	return &examSchedulesHandlers{
		cfg:             cfg,
		examschedulesUS: examschedulesUS,
	}
}

func (e *examSchedulesHandlers) CreateNewExamSchedules() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		examSchedules := &models.ExamSchedules{}

		if err := ctx.ShouldBindJSON(examSchedules); err != nil {
			log.Info().Msg(err.Error())
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := utils.ValidateStruct(ctx, examSchedules)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		res, err := e.examschedulesUS.CreateNewExamSchedules(ctx, examSchedules)
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

func (e *examSchedulesHandlers) FindExamSchedules() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		res, err := e.examschedulesUS.FindExamSchedules(ctx)
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

func (e *examSchedulesHandlers) UpdateExamSchedulesByUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		examSchedules := &models.ExamSchedules{}

		if err := ctx.ShouldBindJSON(examSchedules); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := utils.ValidateStruct(ctx, examSchedules)

		if err != nil {
			errors := utils.ShowErrors(err)
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}

		filter := bson.M{"username": examSchedules.Username}

		res, err := e.examschedulesUS.UpdateExamSchedulesByUsername(ctx, filter, examSchedules)
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
