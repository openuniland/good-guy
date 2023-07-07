package usecase

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/types"
	"github.com/rs/zerolog/log"
)

const loginUrl = "/login.aspx"

type CtmsUS struct {
	cfg *configs.Configs
}

func NewCtmsUseCase(cfg *configs.Configs) ctms.UseCase {
	return &CtmsUS{cfg: cfg}
}

func (us *CtmsUS) Login(ctx context.Context, user *types.LoginRequest) (*types.LoginResponse, error) {

	ctmsUrl := us.cfg.UrlCrawlerList.CtmsUrl

	hash := md5.Sum([]byte(user.Password))
	hashString := hex.EncodeToString(hash[:])

	data := url.Values{
		"__EVENTTARGET":                        {""},
		"__EVENTARGUMENT":                      {""},
		"__VIEWSTATE":                          {"/wEPDwUJNjgxODI3MDEzZGQYhImpueCRmFchkTJkEoLggX4C6Nz/NXMIzR9/49O/0g=="},
		"__VIEWSTATEGENERATOR":                 {"C2EE9ABB"},
		"__EVENTVALIDATION":                    {"/wEdAAQxNFjzuCTBmG4Ry6gmDFTXMVDm8KVzqxEfMx7263Qx5VsdkPb56sD60m4bRwV1zT7o396vFnxqy4G+sdDoX0RYcT0vDsg4dG9gkFX2SUYDeTjkkBvsNMeyzTehazPIVNk="},
		"ctl00$LeftCol$UserLogin1$txtUsername": {user.Username},
		"ctl00$LeftCol$UserLogin1$txtPassword": {hashString},
		"ctl00$LeftCol$UserLogin1$btnLogin":    {"Đăng+nhập"},
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", ctmsUrl+loginUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("error create request")
		return nil, err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,vi;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", ctmsUrl)
	req.Header.Add("Referer", ctmsUrl+loginUrl)
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error send request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error read response")
		return nil, err
	}

	cookie := resp.Header.Get("Set-Cookie")

	if bytes.Contains(body, []byte("Xin chào mừng")) {
		return &types.LoginResponse{
			Cookie:    cookie,
			IsSuccess: true,
		}, nil
	}

	if bytes.Contains(body, []byte("Sai Tên đăng nhập hoặc Mật khẩu")) {
		return &types.LoginResponse{
			Cookie:    "",
			IsSuccess: false,
		}, nil
	}

	return &types.LoginResponse{
		Cookie:    "",
		IsSuccess: false,
	}, nil

}

func (us *CtmsUS) Logout(ctx context.Context, cookie string) error {
	ctmsUrl := us.cfg.UrlCrawlerList.CtmsUrl

	data := url.Values{
		"__VIEWSTATE":          {"/wEPDwUJNjgxODI3MDEzZGQYhImpueCRmFchkTJkEoLggX4C6Nz/NXMIzR9/49O/0g=="},
		"__VIEWSTATEGENERATOR": {"C2EE9ABB"},
		"__CALLBACKID":         {"ctl00$QuanlyMenu1"},
		"__CALLBACKPARAM":      {"logout"},
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", ctmsUrl+loginUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	parts := strings.Split(cookie, ";")
	cookie = strings.TrimSpace(parts[0])

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookie)

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	log.Info().Msg("logout success")

	return nil
}
