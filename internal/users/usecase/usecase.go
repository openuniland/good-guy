package usecase

import (
	"context"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUS struct {
	cfg      *configs.Configs
	userRepo users.Repository
}

func NewUserUseCase(cfg *configs.Configs, userRepo users.Repository) users.UseCase {
	return &UserUS{cfg: cfg, userRepo: userRepo}
}

func (u *UserUS) CreateNewUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	res, err := u.userRepo.Create(ctx, user)

	if err != nil {
		log.Error().Err(err).Msg("error creating new user")
		return nil, err
	}

	return res, nil
}

func (u *UserUS) GetUsers(ctx context.Context) ([]*models.User, error) {
	users, err := u.userRepo.GetUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error getting users")
		return nil, err
	}

	return users, nil
}

func (u *UserUS) GetUserBySubscribedId(ctx context.Context, subscribedId string) (*models.User, error) {
	filter := bson.M{"subscribed_id": subscribedId}
	user, err := u.userRepo.FindOneUserByCondition(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("error getting user by subscribed id")
		return nil, err
	}

	return user, nil
}
