package jobs

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/rs/zerolog/log"
)

func (j *Jobs) getUpcomingExamSchedule() {
	users, err := j.userUC.GetUsers(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("SyncArticles")
		return
	}

	for _, user := range users {
		go func(user *models.User) {

			u := &types.LoginRequest{
				Username: user.Username,
				Password: user.Password,
			}
			err := j.ctmsUS.SendChangedExamScheduleAndNewExamScheduleToUser(context.Background(), u, user.SubscribedID)
			if err != nil {
				log.Error().Err(err).Msg("[JOBS - getUpcomingExamSchedule]" + " - " + user.Username)
				return
			}

			log.Info().Msg("[JOBS - getUpcomingExamSchedule]" + " - " + user.Username)

		}(user)
	}
}
