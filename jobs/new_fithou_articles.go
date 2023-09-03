package jobs

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
	"github.com/rs/zerolog/log"
)

func (j *Jobs) SyncArticles() {
	log.Info().Msg("start SyncArticles")
	res, err := j.articleUC.UpdatedWithNewArticle(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("SyncArticles")
		return
	}

	if !res.IsNew {
		log.Info().Msg("No new articles")
		return
	}

	users, err := j.userUC.GetUsers(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("SyncArticles")
		return
	}

	for _, user := range users {
		go func(user *models.User) {
			for _, article := range res.Data {
				log.Info().Msg("Send message to user: " + user.SubscribedID)

				link := j.cfg.UrlCrawlerList.FithouUrl + article.Link

				message := "📰 " + article.Title + "\n\n" + link + "\n\n"
				err := j.facebookUC.SendTextMessage(context.Background(), user.SubscribedID, message)
				if err != nil {
					log.Error().Err(err).Msg("SyncArticles")
					return
				}
			}

		}(user)
	}
}
