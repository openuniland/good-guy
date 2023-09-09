package jobs

import (
	"context"

	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (j *Jobs) getUpcomingExamSchedule() {
	users, err := j.userUC.GetUsers(context.Background())
	if err != nil {
		log.Error().Err(err).Msgf("[JOBS][ERROR]:[getUpcomingExamSchedule]:[error while getting all users]:[ERROR_INFO=%v]", err)
		return
	}

	for _, user := range users {

		if user.IsExamDay {
			go func(user *models.User) {

				u := &types.LoginCtmsRequest{
					Username: user.Username,
					Password: user.Password,
				}
				err := j.ctmsUS.SendChangedExamScheduleAndNewExamScheduleToUser(context.Background(), u, user.SubscribedID)
				if err != nil {
					log.Error().Err(err).Msgf("[JOBS][ERROR]:[getUpcomingExamSchedule]:[send changed exam schedule]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)

					if err.Error() == constants.NEED_TO_BUY_CTMS && !user.IsDisabled {

						go func() {
							filter := bson.M{"username": user.Username}
							update := bson.M{"is_disabled": true}
							_, err := j.userUC.FindOneAndUpdateUser(context.Background(), filter, update)
							if err != nil {
								log.Error().Err(err).Msgf("[JOBS][ERROR]:[getUpcomingExamSchedule]:[update user]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
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
						log.Error().Err(err).Msgf("[JOBS][ERROR]:[getUpcomingExamSchedule]:[update user is_disabled=false]:[INFO=%s]:[ERROR_INFO=%v]", user.Username, err)
					}
				}
				log.Info().Msgf("[JOBS][INFO]:[getUpcomingExamSchedule]:[INFO=%s]:[SUCCESS]", user.Username)

			}(user)
		}

	}
}
