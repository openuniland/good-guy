package usecase

import (
	"context"
	"fmt"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/auth"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthUS struct {
	cfg        *configs.Configs
	ctmsUC     ctms.UseCase
	userUC     users.UseCase
	facebookUC facebook.UseCase
}

func NewAuthUseCase(cfg *configs.Configs, ctmsUC ctms.UseCase, userUC users.UseCase, facebookUC facebook.UseCase) auth.UseCase {
	return &AuthUS{cfg: cfg, ctmsUC: ctmsUC, userUC: userUC, facebookUC: facebookUC}
}

func (a *AuthUS) Login(ctx context.Context, loginRequest *models.LoginRequest, id string) error {
	user := &types.LoginCtmsRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	}
	res, err := a.ctmsUC.LoginCtms(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("error login ctms")
		return err
	}

	if res.Cookie == "" {
		return fmt.Errorf("error login ctms")
	}

	existedUser, err := a.userUC.GetUserByUsername(ctx, loginRequest.Username)

	if err != nil {
		log.Error().Err(err).Msg("error getting user by username")
		return err
	}

	if existedUser == nil {
		newUser := &models.User{
			Username:     loginRequest.Username,
			Password:     loginRequest.Password,
			SubscribedID: id,
		}

		_, err := a.userUC.CreateNewUser(ctx, newUser)

		if err != nil {
			log.Error().Err(err).Msg("[User not exist]error creating new user")
			return err
		}
	}

	if existedUser != nil {
		if existedUser.SubscribedID != id {
			a.facebookUC.SendTextMessage(ctx, existedUser.SubscribedID, "CTMS BOT: Tài khoản này đã được đăng ký với người dùng khác. Bot sẽ hủy đăng ký tài khoản này.")

		}

		updateUser := &models.User{
			Username:     loginRequest.Username,
			Password:     loginRequest.Password,
			SubscribedID: id,
		}

		filter := bson.M{"username": loginRequest.Username}

		_, err := a.userUC.FindOneAndUpdateUser(ctx, filter, updateUser)
		if err != nil {
			log.Error().Err(err).Msg("[Update user] error updating user")
			return err
		}

	}

	a.facebookUC.SendTextMessage(ctx, id, "CTMS BOT: Đăng nhập thành công.")

	return nil
}
