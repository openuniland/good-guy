package jobs

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Error getting users")
		return
	}

	for _, user := range users {
		if user.IsTrackTimetable {
			go func(user *models.User) {
				u := &types.LoginCtmsRequest{
					Username: user.Username,
					Password: user.Password,
				}

				res, err := j.ctmsUS.LoginCtms(context.Background(), u)
				if err != nil {
					log.Error().Err(err).Msg("Error logging in CTMS")
					return
				}

				dailySchedule, err := j.ctmsUS.GetDailySchedule(context.Background(), res.Cookie)
				if err != nil {
					log.Error().Err(err).Msg("Error getting daily schedule")
					return
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
		log.Error().Err(err).Msg("Error getting users")
		return
	}

	for _, user := range users {
		if user.IsTrackTimetable {
			go func(user *models.User) {
				u := &types.LoginCtmsRequest{
					Username: user.Username,
					Password: user.Password,
				}

				res, err := j.ctmsUS.LoginCtms(context.Background(), u)
				if err != nil {
					log.Error().Err(err).Msg("Error logging in CTMS")
					return
				}

				dailySchedule, err := j.ctmsUS.GetDailySchedule(context.Background(), res.Cookie)
				if err != nil {
					log.Error().Err(err).Msg("Error getting daily schedule")
					return
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
		log.Error().Err(err).Msg("Error getting users")
		return
	}

	for _, user := range users {
		if user.IsTrackTimetable {
			go func(user *models.User) {
				u := &types.LoginCtmsRequest{
					Username: user.Username,
					Password: user.Password,
				}

				res, err := j.ctmsUS.LoginCtms(context.Background(), u)
				if err != nil {
					log.Error().Err(err).Msg("Error logging in CTMS")
					return
				}

				dailySchedule, err := j.ctmsUS.GetDailySchedule(context.Background(), res.Cookie)
				if err != nil {
					log.Error().Err(err).Msg("Error getting daily schedule")
					return
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
