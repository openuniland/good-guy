package jobs

import (
	"context"

	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/rs/zerolog/log"
)

func (j *Jobs) checkCredit() {
	users, err := j.userUC.GetVip(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("[JOBS][ERROR]:[checkCredit]:[get vip users]")
		return
	}

	for _, user := range users {

		go func(user *models.User) {

			u := &types.LoginHouRequest{
				Username:     user.Username,
				Password:     user.Password,
				SubscribedID: user.SubscribedID,
			}
			res, err := j.houUC.LoginHou(context.Background(), u)
			if err != nil {
				log.Error().Err(err).Msg("[JOBS][ERROR]:[checkCredit]:[login hou]")
				return
			}

			if res.SessionId == "" {
				log.Error().Err(err).Msg("[JOBS][ERROR]:[checkCredit]:[session id is empty]")
				return
			}

			credit, err := j.houUC.CheckCreditHou(context.Background(), res.SessionId)
			if err != nil {
				log.Error().Err(err).Msg("[JOBS][ERROR]:[checkCredit]:[check credit hou]")
				return
			}

			if credit == constants.CHANGES_HAVE_OCCURRED {
				j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, "Hình như có đăng ký tín chỉ kìa ku!")
			}

			log.Info().Msgf("[JOBS][SUCCESS]:[checkCredit]:[user=%v]:[credit=%v]", user.Username, credit)
		}(user)

	}
}
