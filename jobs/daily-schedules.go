package jobs

import (
	"context"

	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MORNING   = "07:30 ->11:15"
	AFTERNOON = "13:00 ->16:45"
	EVENING   = "17:15 ->20:00"
)

const (
	MORNING_TEXT   = "Buổi sáng"
	AFTERNOON_TEXT = "Buổi chiều"
	EVENING_TEXT   = "Buổi tối"
)

func (j *Jobs) morningClassSchedule() {
	users, err := j.userUC.GetUsers(context.Background())
	if err != nil {
		log.Error().Err(err).Msgf("[JOBS][ERROR]:[morningClassSchedule]:[error while getting all users]:[ERROR_INFO=%v]", err)
		return
	}

	for _, user := range users {
		if user.IsTrackTimetable {
			go func(user *models.User) {
				u := &types.LoginCtmsRequest{
					Username:     user.Username,
					Password:     user.Password,
					SubscribedID: user.SubscribedID,
				}

				res, err := j.ctmsUS.LoginCtms(context.Background(), u)
				if err != nil {
					log.Error().Err(err).Msgf("[JOBS][ERROR]:[morningClassSchedule]:[error while login ctms]:[ERROR_INFO=%v]", err)
					return
				}

				dailySchedule, err := j.ctmsUS.GetDailySchedule(context.Background(), res.Cookie)
				if err != nil {
					log.Error().Err(err).Msgf("[JOBS][ERROR]:[morningClassSchedule]:[error while getting daily schedule]:[ERROR_INFO=%v]", err)

					if err.Error() == constants.NEED_TO_BUY_CTMS && !user.IsDisabled {

						go func() {
							filter := bson.M{"username": user.Username}
							update := bson.M{"is_disabled": true}
							_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
							if err != nil {
								log.Error().Err(err).Msgf("[JOBS][ERROR]:[morningClassSchedule]:[update user]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
							}
						}()

						go func() {
							j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, constants.NEED_TO_BUY_CTMS_RESPONSE)
						}()
					}
					return
				}

				if user.IsDisabled {
					filter := bson.M{"username": user.Username}
					update := bson.M{"is_disabled": false}
					_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
					if err != nil {
						log.Error().Err(err).Msgf("[JOBS][ERROR]:[morningClassSchedule]:[update user is_disabled=false]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
					}
				}

				for _, schedule := range dailySchedule {
					if schedule.Time == MORNING {
						go j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, utils.DailyScheduleMessage(utils.GetClassStatus(schedule.Status, MORNING_TEXT), schedule))
					}
				}

			}(user)
		}
	}
}

func (j *Jobs) afternoonClassSchedule() {
	users, err := j.userUC.GetUsers(context.Background())
	if err != nil {
		log.Error().Err(err).Msgf("[JOBS][ERROR]:[afternoonClassSchedule]:[error while getting all users]:[ERROR_INFO=%v]", err)
		return
	}

	for _, user := range users {
		if user.IsTrackTimetable {
			go func(user *models.User) {
				u := &types.LoginCtmsRequest{
					Username:     user.Username,
					Password:     user.Password,
					SubscribedID: user.SubscribedID,
				}

				res, err := j.ctmsUS.LoginCtms(context.Background(), u)
				if err != nil {
					log.Error().Err(err).Msg("Error logging in CTMS")
					return
				}

				dailySchedule, err := j.ctmsUS.GetDailySchedule(context.Background(), res.Cookie)
				if err != nil {
					log.Error().Err(err).Msg("Error getting daily schedule")

					if err.Error() == constants.NEED_TO_BUY_CTMS && !user.IsDisabled {

						go func() {
							filter := bson.M{"username": user.Username}
							update := bson.M{"is_disabled": true}
							_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
							if err != nil {
								log.Error().Err(err).Msgf("[JOBS][ERROR]:[afternoonClassSchedule]:[update user]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
							}
						}()

						go func() {
							j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, constants.NEED_TO_BUY_CTMS_RESPONSE)
						}()
					}
					return
				}

				if user.IsDisabled {
					filter := bson.M{"username": user.Username}
					update := bson.M{"is_disabled": false}
					_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
					if err != nil {
						log.Error().Err(err).Msgf("[JOBS][ERROR]:[afternoonClassSchedule]:[update user is_disabled=false]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
					}
				}

				for _, schedule := range dailySchedule {
					if schedule.Time == AFTERNOON {
						go j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, utils.DailyScheduleMessage(utils.GetClassStatus(schedule.Status, AFTERNOON_TEXT), schedule))
					}
				}

			}(user)
		}
	}
}

func (j *Jobs) eveningClassSchedule() {
	users, err := j.userUC.GetUsers(context.Background())
	if err != nil {
		log.Error().Err(err).Msgf("[JOBS][ERROR]:[eveningClassSchedule]:[error while getting all users]:[ERROR_INFO=%v]", err)
		return
	}

	for _, user := range users {
		if user.IsTrackTimetable {
			go func(user *models.User) {
				u := &types.LoginCtmsRequest{
					Username:     user.Username,
					Password:     user.Password,
					SubscribedID: user.SubscribedID,
				}

				res, err := j.ctmsUS.LoginCtms(context.Background(), u)
				if err != nil {
					log.Error().Err(err).Msg("Error logging in CTMS")
					return
				}

				dailySchedule, err := j.ctmsUS.GetDailySchedule(context.Background(), res.Cookie)
				if err != nil {
					log.Error().Err(err).Msg("Error getting daily schedule")

					if err.Error() == constants.NEED_TO_BUY_CTMS && !user.IsDisabled {

						go func() {
							filter := bson.M{"username": user.Username}
							update := bson.M{"is_disabled": true}
							_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
							if err != nil {
								log.Error().Err(err).Msgf("[JOBS][ERROR]:[afternoonClassSchedule]:[update user]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
							}
						}()

						go func() {
							j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, constants.NEED_TO_BUY_CTMS_RESPONSE)
						}()
					}
					return
				}

				if user.IsDisabled {
					filter := bson.M{"username": user.Username}
					update := bson.M{"is_disabled": false}
					_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
					if err != nil {
						log.Error().Err(err).Msgf("[JOBS][ERROR]:[afternoonClassSchedule]:[update user is_disabled=false]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
					}
				}

				for _, schedule := range dailySchedule {
					if schedule.Time == EVENING {
						go j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, utils.DailyScheduleMessage(utils.GetClassStatus(schedule.Status, EVENING_TEXT), schedule))
					}
				}

			}(user)
		}
	}
}
