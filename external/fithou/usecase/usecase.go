package usecase

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/fithou"
	"github.com/openuniland/good-guy/external/types"
	"github.com/rs/zerolog/log"
)

type FithouUS struct {
	cfg *configs.Configs
}

func NewFithouUseCase(cfg *configs.Configs) fithou.UseCase {
	return &FithouUS{cfg: cfg}
}

func (us *FithouUS) CrawlArticlesFromFirstPage(ctx context.Context) ([]*types.ArticleCrawl, error) {
	res, err := http.Get(us.cfg.UrlCrawlerList.FithouCategoriesUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error().Msgf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error().Msgf("goquery error: %s", err.Error())
		return nil, err
	}

	var articles []*types.ArticleCrawl

	doc.Find("#LeftCol_pnlCategory div[class=article]").Each(func(_ int, s *goquery.Selection) {
		title := s.Find("a").Text()
		link, _ := s.Find("a").Attr("href")
		re := regexp.MustCompile(`aid=(\d+)`)
		aidArray := re.FindStringSubmatch(link)

		aid, err := strconv.Atoi(aidArray[1])
		if err != nil {
			log.Error().Msgf("convert string to int error: %s", err.Error())
		}

		articles = append(articles, &types.ArticleCrawl{
			Title: strings.TrimSpace(title),
			Link:  link,
			Aid:   aid,
		})
	})

	return articles, nil
}
