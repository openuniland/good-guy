package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/hou"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type HouUS struct {
	cfg      *configs.Configs
	userRepo users.Repository
}

func NewHouUseCase(cfg *configs.Configs, userRepo users.Repository) hou.UseCase {
	return &HouUS{cfg: cfg, userRepo: userRepo}
}

func (h *HouUS) LoginHou(ctx context.Context, user *types.LoginHouRequest) (*types.LoginHouResponse, error) {
	// Create a cookie jar to store cookies
	cookieJar, _ := cookiejar.New(nil)

	// Create a client with cookie jar
	client := &http.Client{
		Jar: cookieJar,
	}

	targetURL := fmt.Sprintf("%s%s", h.cfg.UrlCrawlerList.SinhVienUrl, constants.LOGIN)

	resp, err := client.Get(targetURL)
	if err != nil {
		log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[client.Get]:[ERROR_INFO=%v]", err)
		return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[goquery]:[ERROR_INFO=%v]", err)
		return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
	}

	// find input with name="execution" type="hidden"
	execution := doc.Find("input[name='execution']").AttrOr("value", "")
	log.Info().Msgf("[INFO]:[USECASE]:[LoginHou]:[execution]:[INFO=%v]", execution)
	defer resp.Body.Close()

	if execution == "" {
		log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[execution]:[ERROR_INFO=%s]", "execution is empty")
		return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
	}

	response := &types.LoginHouResponse{
		Username:  user.Username,
		SessionId: "",
		AspxAuth:  "",
	}

	// check if redirected
	if resp.Request.URL.String() != targetURL {

		loginURL := fmt.Sprintf("%s/cas/login?service=%s", h.cfg.UrlCrawlerList.CasHouUrl, targetURL)
		resp, err := client.PostForm(loginURL, url.Values{
			"username":  {user.Username},
			"password":  {user.Password},
			"execution": {execution},
			"_eventId":  {"submit"},
		})
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[client.PostForm]:[ERROR_INFO=%v, DATA=%v]", err, user)
			return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[goquery]:[ERROR_INFO=%v]", err)
			return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
		}

		loginUnsuccessfulText := doc.Find("div.alert.alert-danger span").Text()
		if strings.TrimSpace(loginUnsuccessfulText) == constants.LOGIN_UNSUCCESSFUL {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[loginUnsuccessfulText]:[DATA=%v]", user)
			return nil, errors.New(constants.INCORRECCT_USERNAME_OR_PASSWORD)
		}

		log.Info().Msgf("[INFO]:[USECASE]:[LoginHou]:[client.PostForm]:[DATA=%v]", loginUnsuccessfulText)

		defer resp.Body.Close()

		// read cookies
		cookies := cookieJar.Cookies(resp.Request.URL)
		for _, cookie := range cookies {
			fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
			if cookie.Name == constants.ASPNET_SESSION_ID {
				response.SessionId = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
			}

			if cookie.Name == constants.ASPXAUTH {
				response.AspxAuth = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
			}
		}

		if response.SessionId == "" || response.AspxAuth == "" {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[response.SessionId == '']:[ERROR_INFO=%v, DATA=%v]", err, user)
			return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
		}

		filter := bson.M{"subscribed_id": user.SubscribedID}
		userRecord, err := h.userRepo.FindOneUserByCondition(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[userRepo.FindOneUserByCondition]:[ERROR_INFO=%v, DATA=%v]", err, user)
			return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
		}

		update := bson.M{"session_id": response.SessionId, "aspx_auth": response.AspxAuth}
		if user.Username != userRecord.Username {
			update = bson.M{"session_id": response.SessionId, "aspx_auth": response.AspxAuth, "username": user.Username, "password": user.Password, "login_provider": constants.SINHVIEN}
		}

		user, err := h.userRepo.FindOneAndUpdate(ctx, filter, update)
		if err != nil {
			log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[userRepo.FindOneAndUpdate]:[ERROR_INFO=%v, DATA=%v]", err, user)
			return nil, errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
		}

		log.Info().Msgf("[INFO]:[USECASE]:[LoginHou]:[userRepo.FindOneAndUpdate]:[INFO=%v]", user)
	} else {
		log.Warn().Msgf("[WARN]:[USECASE]:[LoginHou]:[check if redirected]:[WARN_INFO=%s]", "There are no redirects.")
	}

	return response, nil
}

func (h *HouUS) LogoutHou(ctx context.Context, SessionId string) error {
	client := &http.Client{}

	targetURL := fmt.Sprintf("%s%s", h.cfg.UrlCrawlerList.SinhVienUrl, constants.LOGOUT)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {

		return errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
	}
	// Set request headers
	req.Header.Set("Cookie", SessionId)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(constants.LOGIN_UNSUCCESSFUL_ALERT)
	}
	defer resp.Body.Close()

	return nil
}
