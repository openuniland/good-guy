package usecase

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/hou"
	"github.com/openuniland/good-guy/external/types"
	"github.com/rs/zerolog/log"
)

type HouUS struct {
	cfg *configs.Configs
}

func NewHouUseCase(cfg *configs.Configs) hou.UseCase {
	return &HouUS{cfg: cfg}
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
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[goquery]:[ERROR_INFO=%v]", err)
		return nil, err
	}

	// find input with name="execution" type="hidden"
	execution := doc.Find("input[name='execution']").AttrOr("value", "")
	log.Info().Msgf("[INFO]:[USECASE]:[LoginHou]:[execution]:[INFO=%v]", execution)
	defer resp.Body.Close()

	if execution == "" {
		log.Error().Err(err).Msgf("[ERROR]:[USECASE]:[LoginHou]:[execution]:[ERROR_INFO=%s]", "execution is empty")
		return nil, nil
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
			return nil, err
		}
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
	} else {
		log.Warn().Msgf("[WARN]:[USECASE]:[LoginHou]:[check if redirected]:[WARN_INFO=%s]", "There are no redirects.")
	}

	return response, nil
}

func (h *HouUS) LogoutHou(ctx context.Context) error {
	return nil
}
