package jobs

import (
	"context"
	"fmt"

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

			var u types.LoginRequest
			u.Username = user.Username
			u.Password = user.Password

			upcoming, err := j.ctmsUS.GetUpcomingExamSchedule(context.Background(), &u)
			if err != nil {
				log.Error().Err(err).Msg("SyncArticles")
				return
			}
			fmt.Println(upcoming)

			// if len(upcoming) == 0 {
			// 	return
			// }

			message := "	"
			err = j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, message)
			if err != nil {
				log.Error().Err(err).Msg("SyncArticles")
				return
			}

		}(user)
	}
}
