package jobs

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (j *Jobs) syncArticles() {
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

	for _, id := range res.SubscribedIDs {

		go func(id string) {
			for _, article := range res.Data {
				log.Info().Msg("Send message to user: " + id)

				link := j.cfg.UrlCrawlerList.FithouUrl + article.Link

				message := "📰 " + article.Title + "\n\n" + link + "\n\n"
				err := j.facebookUC.SendTextMessage(context.Background(), id, message)
				if err != nil {
					log.Error().Err(err).Msg("SyncArticles")
					return
				}
			}

		}(id)
	}
}
