package usecase

import (
	"context"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/cookies"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CookieUS struct {
	cfg        *configs.Configs
	cookieRepo cookies.Repository
}

func NewCookieUseCase(cfg *configs.Configs, cookieRepo cookies.Repository) cookies.UseCase {
	return &CookieUS{cfg: cfg, cookieRepo: cookieRepo}
}

func (c *CookieUS) CreateNewCookie(ctx context.Context, cookie *models.Cookie) (*mongo.InsertOneResult, error) {

	_, err := c.cookieRepo.Create(ctx, cookie)
	if err != nil {
		log.Err(err).Msgf("[ERROR]:[USECASE]:[CreateNewCookie]:[INFO=%v]:[ERROR_INFO=%v]", cookie.Username, err)
		return nil, err
	}

	log.Info().Msgf("[INFO]:[USECASE]:[CreateNewCookie]:[INFO=%v]:[SUCCESS]", cookie.Username)
	return nil, nil
}

func (c *CookieUS) FindOneCookie(ctx context.Context, username string) (*models.Cookie, error) {

	filter := bson.M{"username": username}

	cookie, err := c.cookieRepo.FindOne(ctx, filter)
	if err != nil {
		log.Err(err).Msgf("[ERROR]:[USECASE]:[FindOneCookie]:[INFO=%v]:[ERROR_INFO=%v]", filter, err)
		return nil, err
	}

	log.Info().Msgf("[INFO]:[USECASE]:[FindOneCookie]:[INFO=%v]:[SUCCESS]", filter)
	return cookie, nil
}

func (c *CookieUS) UpdateSertCookie(ctx context.Context, filter bson.M, update bson.M) error {

	_, err := c.cookieRepo.UpdateSertOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("[ERROR]:[USECASE]:[UpdateCookie]:[FILTER=%v, UPDATE=%v]:[ERROR_INFO=%v]", filter, update, err)
		return err
	}

	log.Info().Msgf("[INFO]:[USECASE]:[UpdateCookie]:[FILTER=%v, UPDATE=%v]:[SUCCESS]", filter, update)
	return nil
}
