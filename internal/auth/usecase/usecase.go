package usecase

import (
	"context"
	"fmt"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/auth"
	"github.com/openuniland/good-guy/internal/cookies"
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
	cookieUC   cookies.UseCase
}

func NewAuthUseCase(cfg *configs.Configs, ctmsUC ctms.UseCase, userUC users.UseCase, facebookUC facebook.UseCase, cookieUC cookies.UseCase) auth.UseCase {
	return &AuthUS{cfg: cfg, ctmsUC: ctmsUC, userUC: userUC, facebookUC: facebookUC, cookieUC: cookieUC}
}

func (a *AuthUS) Login(ctx context.Context, loginRequest *models.LoginRequest) error {
	user := &types.LoginCtmsRequest{
		Username:     loginRequest.Username,
		Password:     loginRequest.Password,
		SubscribedID: loginRequest.Id,
	}
	res, err := a.ctmsUC.LoginCtms(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("error login ctms")
		return err
	}

	go func() {
		cookie := &models.Cookie{
			Username: user.Username,
			Cookies:  []string{res.Cookie},
		}

		_, err = a.cookieUC.CreateNewCookie(ctx, cookie)
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginCtms]:[create new cookie]:[INFO=%s]:[ERROR_INFO%v]", user.Username, err)
			return
		}
	}()

	go func() {
		// [LOGOUT_CTMS]
		err = a.ctmsUC.LogoutCtms(ctx, res.Cookie)
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[Login]:[error while login]:[INFO=%s, COOKIE=%s]:[%v]", user.Username, res.Cookie, err)
		}
		log.Info().Msgf("[INFO]:[Login]:[logout successful]:[INFO=%s, COOKIE=%s]", user.Username, res.Cookie)
	}()

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
			SubscribedID: loginRequest.Id,
		}

		_, err := a.userUC.CreateNewUser(ctx, newUser)

		if err != nil {
			log.Error().Err(err).Msg("[User not exist] error creating new user")
			return err
		}
	}

	if existedUser != nil {
		if existedUser.SubscribedID != loginRequest.Id {
			a.facebookUC.SendTextMessage(ctx, existedUser.SubscribedID, "CTMS BOT: Tài khoản này đã được đăng ký với người dùng khác. Bot sẽ hủy đăng ký tài khoản này.")

		}

		updateUser := bson.M{
			"username":      loginRequest.Username,
			"password":      loginRequest.Password,
			"subscribed_id": loginRequest.Id,
		}

		filter := bson.M{"username": loginRequest.Username}

		_, err := a.userUC.FindOneAndUpdateUser(ctx, filter, updateUser)
		if err != nil {
			log.Error().Err(err).Msg("[Update user] error updating user")
			return err
		}

	}

	a.facebookUC.SendTextMessage(ctx, loginRequest.Id, "CTMS BOT: Đăng nhập thành công.")

	return nil
}
