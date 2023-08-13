package usecase

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/types"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

const loginUrl = "/login.aspx"

const SCHOOL_SCHEDULE_URL = "http://ctms.fithou.net.vn/Lichhoc.aspx?sid="
const EXPIRED_CTMS = "Từ 2/2022, hãy thực hiện theo thông báo này để nhận được sự Hỗ trợ duy trì tài khoản truy cập CTMS từ khoa CNTT."
const SESSION_EXPIRED_MESSAGE = "Phiên làm việc hết hạn hoặc Bạn không có quyền truy cập chức năng này"

type CtmsUS struct {
	cfg             *configs.Configs
	examschedulesUS examschedules.UseCase
}

func NewCtmsUseCase(cfg *configs.Configs, examschedulesUS examschedules.UseCase) ctms.UseCase {
	return &CtmsUS{cfg: cfg, examschedulesUS: examschedulesUS}
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
		log.Error().Msg("error create request login" + err.Error())
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
		log.Err(err).Msg("error send request login" + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("error read body login" + err.Error())
		return nil, err
	}

	cookie := resp.Header.Get("Set-Cookie")

	if bytes.Contains(body, []byte("Xin chào mừng")) {
		return &types.LoginResponse{
			Cookie: cookie,
		}, nil
	}

	if bytes.Contains(body, []byte("Sai Tên đăng nhập hoặc Mật khẩu")) {
		return nil, errors.New("wrong username or password")
	}

	return nil, errors.New("an unknown error")

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

		log.Err(err).Msg("error create request to logout" + err.Error())
		return err
	}

	parts := strings.Split(cookie, ";")
	cookie = strings.TrimSpace(parts[0])

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookie)

	_, err = client.Do(req)
	if err != nil {
		log.Err(err).Msg("error send request to logout" + err.Error())
		return err
	}

	log.Info().Msg("logout success")

	return nil
}

func (us *CtmsUS) GetDailySchedule(ctx context.Context, cookie string) ([]*types.DailySchedule, error) {

	date := utils.FormatDateTimeToGetDailySchedule()

	data := url.Values{
		"__EVENTTARGET":                         {""},
		"__EVENTARGUMENT":                       {""},
		"__VIEWSTATE":                           {"/wEPDwUKMTA4NDM3NDc2OGQYBwUzY3RsMDAkTGVmdENvbCRMaWNoaG9jMSRycHRyTGljaGhvYyRjdGwwNSRncnZMaWNoaG9jDzwrAAwBCAIBZAUzY3RsMDAkTGVmdENvbCRMaWNoaG9jMSRycHRyTGljaGhvYyRjdGwwMyRncnZMaWNoaG9jDzwrAAwBCGZkBTNjdGwwMCRMZWZ0Q29sJExpY2hob2MxJHJwdHJMaWNoaG9jJGN0bDAyJGdydkxpY2hob2MPPCsADAEIAgFkBTNjdGwwMCRMZWZ0Q29sJExpY2hob2MxJHJwdHJMaWNoaG9jJGN0bDA2JGdydkxpY2hob2MPPCsADAEIZmQFM2N0bDAwJExlZnRDb2wkTGljaGhvYzEkcnB0ckxpY2hob2MkY3RsMDEkZ3J2TGljaGhvYw88KwAMAQhmZAUzY3RsMDAkTGVmdENvbCRMaWNoaG9jMSRycHRyTGljaGhvYyRjdGwwNCRncnZMaWNoaG9jDzwrAAwBCGZkBTNjdGwwMCRMZWZ0Q29sJExpY2hob2MxJHJwdHJMaWNoaG9jJGN0bDAwJGdydkxpY2hob2MPPCsADAEIAgFkhO4CQTCT9FOotSw2ZoTf5gEBbXaed4Q4OAV5jtaoJYE="},
		"__VIEWSTATEGENERATOR":                  {"CB78C13A"},
		"__EVENTVALIDATION":                     {"/wEdAAPwrTvSkjO6MxCyG5nv8RpLWWWHEhzFiGyQmAroNHRecPGp81KLC9U2/agHpgpfb4atL2GQMaATghjy+bylAXhJAkV++jXlveksbno26k3dtg=="},
		"ctl00$LeftCol$Lichhoc1$txtNgaydautuan": {date},
		"ctl00$LeftCol$Lichhoc1$btnXemlich":     {"Xem+lịch"},
	}

	// Create HTTP client
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest("POST", SCHOOL_SCHEDULE_URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Err(err).Msg("error create request to get daily schedule" + err.Error())
		return nil, err
	}
	// Set request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Origin", us.cfg.UrlCrawlerList.CtmsUrl)
	req.Header.Set("Referer", SCHOOL_SCHEDULE_URL)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("error send request to get daily schedule" + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Err(err).Msg("error parse response to get daily schedule" + err.Error())
		return nil, err
	}

	NoPermissionText := doc.Find(".NoPermission h3").Text()
	if strings.TrimSpace(NoPermissionText) == SESSION_EXPIRED_MESSAGE {

		log.Error().Msg("session expired")
		return nil, errors.New("session expired")
	}

	expiredNotiText := doc.Find("#leftcontent #thongbao").Text()
	if strings.TrimSpace(expiredNotiText) == EXPIRED_CTMS {

		log.Error().Msg("need to buy ctm")
		return nil, errors.New("need to buy ctm")
	}

	var dailyScheduleData []*types.DailySchedule
	doc.Find("#leftcontent #LeftCol_Lichhoc1_pnView").ChildrenFiltered("div").Each(func(_ int, s *goquery.Selection) {
		day := s.First().Find("b").Text()

		today := utils.TodayFormatted()

		words := strings.Split(day, "\n")
		date := ""

		if len(words) >= 2 {
			date = strings.TrimSpace(words[2])
		}

		if today == date {
			s.Find("div table tbody tr").Each(func(j int, ss *goquery.Selection) {

				if j != 0 {
					res := &types.DailySchedule{
						SerialNumber: strings.TrimSpace(ss.Find("td").Eq(0).Text()),
						Time:         strings.TrimSpace(ss.Find("td").Eq(1).Text()),
						ClassRoom:    strings.TrimSpace(ss.Find("td").Eq(2).Text()),
						SubjectName:  strings.TrimSpace(ss.Find("td").Eq(3).Text()),
						Lecturer:     strings.TrimSpace(ss.Find("td").Eq(4).Text()),
						ClassCode:    strings.TrimSpace(ss.Find("td").Eq(5).Text()),
						Status:       strings.TrimSpace(ss.Find("td").Eq(6).Text()),
					}

					dailyScheduleData = append(dailyScheduleData, res)
				}

			})
		}

	})

	return dailyScheduleData, nil
}

func (us *CtmsUS) GetExamSchedule(ctx context.Context, cookie string) ([]*types.ExamSchedule, error) {

	examScheduleUrl := us.cfg.UrlCrawlerList.ExamScheduleUrl

	// Create HTTP client
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest("GET", examScheduleUrl, nil)
	if err != nil {
		log.Err(err).Msg("error create request to get exam schedule" + err.Error())
		return nil, err
	}
	// Set request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Origin", examScheduleUrl)
	req.Header.Set("Referer", examScheduleUrl)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	fmt.Println("examScheduleUrl", examScheduleUrl)
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("error send request to get exam schedule: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Err(err).Msg("error parse response to get exam schedule" + err.Error())
		return nil, err
	}

	NoPermissionText := doc.Find(".NoPermission h3").Text()
	if strings.TrimSpace(NoPermissionText) == SESSION_EXPIRED_MESSAGE {

		log.Error().Msg("session expired")
		return nil, errors.New("session expired")
	}

	expiredNotiText := doc.Find("#leftcontent #thongbao").Text()
	if strings.TrimSpace(expiredNotiText) == EXPIRED_CTMS {

		log.Error().Msg("need to buy ctm")
		return nil, errors.New("need to buy ctm")
	}

	var examScheduleData []*types.ExamSchedule
	doc.Find(".RowEffect tbody tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			res := &types.ExamSchedule{
				SerialNumber: strings.TrimSpace(s.Find("td").Eq(0).Text()),
				Time:         strings.TrimSpace(s.Find("td").Eq(1).Text()),
				ClassRoom:    strings.TrimSpace(s.Find("td").Eq(2).Text()),
				SubjectName:  strings.TrimSpace(s.Find("td").Eq(3).Text()),
				ExamListCode: strings.TrimSpace(s.Find("td").Eq(4).Text()),
			}

			examScheduleData = append(examScheduleData, res)
		}
	})

	return examScheduleData, nil
}

func (us *CtmsUS) GetUpcomingExamSchedule(ctx context.Context, user *types.LoginRequest) ([]*types.ExamSchedule, error) {
	cookie, err := us.Login(ctx, user)
	if err != nil {
		log.Err(err).Msg("error login to get upcoming exam schedule")
		return nil, err
	}

	examSchedule, err := us.GetExamSchedule(ctx, cookie.Cookie)
	if err != nil {
		log.Err(err).Msg("error get exam schedule to get upcoming exam schedule")
		return nil, err
	}

	filter := bson.M{"username": user.Username}

	res, err := us.examschedulesUS.FindExamSchedulesByUsername(ctx, filter)
	if err != nil {
		log.Err(err).Msg("error find exam schedule by username to get upcoming exam schedule")
		return nil, err
	}

	fmt.Println("res", res)

	return examSchedule, nil
}
