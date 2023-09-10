package usecase

import (
	"context"
	"errors"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/articles"
	"github.com/openuniland/good-guy/internal/common"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type CommonUS struct {
	cfg        *configs.Configs
	facebookUS facebook.UseCase
	ctmsUC     ctms.UseCase
	userUC     users.UseCase
	articleUC  articles.UseCase
}

func NewCommonUseCase(cfg *configs.Configs, facebookUS facebook.UseCase, ctmsUC ctms.UseCase, userUC users.UseCase, articleUC articles.UseCase) common.UseCase {
	return &CommonUS{cfg: cfg, facebookUS: facebookUS, ctmsUC: ctmsUC, userUC: userUC, articleUC: articleUC}
}

func (us *CommonUS) GetNotificationOfExamSchedule(ctx context.Context, id string) error {
	filter := bson.M{"subscribed_id": id}
	update := bson.M{"is_exam_day": true}

	_, err := us.userUC.FindOneAndUpdateUser(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("error find one and update user")
		us.facebookUS.SendTextMessage(ctx, id, "❗️ Bạn chưa thêm tài khoản CTMS vào hệ thống.")
		return err
	}

	us.facebookUS.SendTextMessage(ctx, id, "🔔 Bật chức năng thông báo lịch thi thành công!")

	return nil
}

func (us *CommonUS) CancelGetNotificationOfExamSchedule(ctx context.Context, id string) error {
	filter := bson.M{"subscribed_id": id}
	update := bson.M{"is_exam_day": false}

	_, err := us.userUC.FindOneAndUpdateUser(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("error find one and update user")
		us.facebookUS.SendTextMessage(ctx, id, "❗️ Bạn chưa thêm tài khoản CTMS vào hệ thống.")
		return err
	}

	us.facebookUS.SendTextMessage(ctx, id, "🔕 Đã tắt chức năng thông báo lịch thi!")
	return nil
}

func (us *CommonUS) AddCtmsTimetableService(ctx context.Context, id string) error {

	filter := bson.M{"subscribed_id": id}
	update := bson.M{"is_track_timetable": true}

	_, err := us.userUC.FindOneAndUpdateUser(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("error find one and update user")
		us.facebookUS.SendTextMessage(ctx, id, "❗️ Bạn chưa thêm tài khoản CTMS vào hệ thống.")
		return err
	}

	us.facebookUS.SendTextMessage(ctx, id, "🔔 Bật chức năng thông báo lịch học hàng ngày thành công!")
	return nil

}

func (us *CommonUS) RemoveCtmsTimetableService(ctx context.Context, id string) error {

	filter := bson.M{"subscribed_id": id}
	update := bson.M{"is_track_timetable": false}

	_, err := us.userUC.FindOneAndUpdateUser(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("error find one and update user")
		us.facebookUS.SendTextMessage(ctx, id, "❗️ Bạn chưa thêm tài khoản CTMS vào hệ thống.")
		return err
	}

	us.facebookUS.SendTextMessage(ctx, id, "🔕 Đã tắt chức năng thông báo lịch học hàng ngày!")
	return nil

}

func (us *CommonUS) AddFithouCrawlService(ctx context.Context, id string) error {
	article, err := us.articleUC.FindOne(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error find one article")
		return err
	}

	for _, item := range article.SubscribedIDs {
		if item == id {
			us.facebookUS.SendTextMessage(ctx, id, "Bạn đã đăng ký nhận thông báo từ FITHOU rồi!")
			return nil
		}
	}

	err = us.articleUC.AddNewSubscriber(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("error add new subscriber")
		return err
	}

	go us.facebookUS.SendTextMessage(ctx, id, "Đăng ký nhận thông báo từ FITHOU thành công!")

	link := us.cfg.UrlCrawlerList.FithouUrl + article.Link
	message := "📰 " + article.Title + "\n\n" + link + "\n\n"
	go us.facebookUS.SendTextMessage(ctx, id, "Gửi bạn bài viết mới nhất hiện tại. Bot sẽ câp nhật thông báo khi có bài viết mới."+"\n"+message)

	return nil
}

func (us *CommonUS) RemoveFithouCrawlService(ctx context.Context, id string) error {

	err := us.articleUC.RemoveSubscriber(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("error remove subscriber")
		return err
	}

	us.facebookUS.SendTextMessage(ctx, id, "Hủy nhận thông báo từ FITHOU thành công!")

	return nil
}

func (us *CommonUS) RemoveUser(ctx context.Context, id string) error {
	filter := bson.M{"subscribed_id": id}

	_, err := us.userUC.FindOneAndDeleteUser(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msg("error remove user")
		return err
	}

	us.facebookUS.SendTextMessage(ctx, id, "Đã xóa tài khoản CTMS thành công!")
	log.Info().Msg("[success]" + "-" + "[remove user]" + "-" + "[" + id + "]")
	return nil
}

func (us *CommonUS) SendLoginCtmsButton(ctx context.Context, id string) error {
	user, err := us.userUC.GetUserBySubscribedId(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("error get user by subscribed id")
		return err
	}

	if user != nil {
		us.facebookUS.SendTextMessage(ctx, id, "Bạn đã đăng nhập CTMS rồi!")
		return nil
	}

	input := &types.SendButtonMessageRequest{
		ImageUrl: constants.IMAGE_URL_LOGIN_CTMS_BTN,
		Title:    "Đăng nhập CTMS",
		Url:      us.cfg.Server.Host + "?id=" + id,
		Subtitle: "Đăng nhập để nhận thông báo từ CTMS",
		BtnText:  "Đăng nhập",
	}

	us.facebookUS.SendButtonMessage(ctx, id, input)
	return nil
}

func (us *CommonUS) VerifyFacebookWebhook(ctx context.Context, token, challenge string) (string, error) {
	if token == us.cfg.FBConfig.FBVerifyToken {
		return challenge, nil
	}

	return "", errors.New("error verify token")
}

func (us *CommonUS) ChatScript(ctx context.Context, id string, msg string) {
	switch msg {
	case utils.Command.Login:
		us.SendLoginCtmsButton(ctx, id)
		return
	case utils.Command.Remove:
		us.RemoveUser(ctx, id)
		return
	case utils.Command.Help:
		us.facebookUS.SendTextMessage(ctx, id, utils.HelpScript())
		return
	case utils.Command.Examday:
		us.GetNotificationOfExamSchedule(ctx, id)
		return
	case utils.Command.UnExamday:
		us.CancelGetNotificationOfExamSchedule(ctx, id)
		return
	case utils.Command.ForceLogout:
		//
		return
	default:
		us.facebookUS.SendTextMessage(ctx, id, "Bot ngu ngok quá, không hiểu gì hết :(. "+"\n"+" /help để biết thêm chi tiết!")
		return
	}

}

func (us *CommonUS) HandleFacebookWebhook(ctx context.Context, data *types.FacebookWebhookRequest) error {

	if data.Object != "page" {
		log.Error().Msg("error object is not page")
		return nil
	}

	messaging := data.Entry[0].Messaging

	for _, element := range messaging {
		sender := element.Sender
		postback := element.Postback
		msg := element.Message

		id := sender.ID

		if postback != nil {
			switch postback.Payload {
			case "GET_STARTED":
				us.facebookUS.SendTextMessage(ctx, id, "Chào mừng bạn đến với Fithou BOT. Chúc bạn có một trải nghiệm zui zẻ :D. /help để biết thêm chi tiết!")
				return nil
			case "HELP":
				us.facebookUS.SendTextMessage(ctx, id, utils.HelpScript())
				return nil
			case "CTMS_SERVICE":
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReplyRequest{{
					ContentType: "text",
					Title:       "Thêm tài khoản CTMS",
					Payload:     "ADD_CTMS_ACCOUNT",
					ImageUrl:    constants.NOTI_IMAGE_ON,
				}, {
					ContentType: "text",
					Title:       "Xóa tài khoản CTMS",
					Payload:     "REMOVE_CTMS_ACCOUNT",
					ImageUrl:    constants.NOTI_IMAGE_OFF,
				}})
				return nil
			case "FITHOU_CRAWL_SERVICE":
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReplyRequest{{
					ContentType: "text",
					Title:       "Bật thông báo",
					Payload:     "ADD_FITHOU_CRAWL_SERVICE",
					ImageUrl:    constants.NOTI_IMAGE_ON,
				}, {
					ContentType: "text",
					Title:       "Tắt thông báo",
					Payload:     "REMOVE_FITHOU_CRAWL_SERVICE",
					ImageUrl:    constants.NOTI_IMAGE_OFF,
				}})
				return nil
			case "CTMS_CREDITS_SERVICE":
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReplyRequest{{
					ContentType: "text",
					Title:       "Bật theo dõi",
					Payload:     "ADD_CTMS_CREDITS_SERVICE",
					ImageUrl:    constants.NOTI_IMAGE_ON,
				}, {
					ContentType: "text",
					Title:       "Tắt theo dõi",
					Payload:     "REMOVE_CTMS_CREDITS_SERVICE",
					ImageUrl:    constants.NOTI_IMAGE_OFF,
				}})
				return nil
			case "CTMS_TIMETABLE_SERVICE":
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReplyRequest{{
					ContentType: "text",
					Title:       "Bật thông báo",
					Payload:     "ADD_CTMS_TIMETABLE_SERVICE",
					ImageUrl:    constants.NOTI_IMAGE_ON,
				}, {
					ContentType: "text",
					Title:       "Tắt thông báo",
					Payload:     "REMOVE_CTMS_TIMETABLE_SERVICE",
					ImageUrl:    constants.NOTI_IMAGE_OFF,
				}})
				return nil
			default:
				return nil
			}
		} else if msg != nil {

			quickReply := msg.QuickReply

			if quickReply != nil {
				switch quickReply.Payload {
				case "ADD_CTMS_ACCOUNT":
					us.SendLoginCtmsButton(ctx, id)
					return nil
				case "REMOVE_CTMS_ACCOUNT":
					us.RemoveUser(ctx, id)
					return nil
				case "ADD_FITHOU_CRAWL_SERVICE":
					us.AddFithouCrawlService(ctx, id)
					return nil
				case "REMOVE_FITHOU_CRAWL_SERVICE":
					us.RemoveFithouCrawlService(ctx, id)
					return nil
				case "ADD_CTMS_CREDITS_SERVICE":
					us.facebookUS.SendTextMessage(ctx, id, "Chức năng dành cho quản trị viên!")
					return nil
				case "REMOVE_CTMS_CREDITS_SERVICE":
					us.facebookUS.SendTextMessage(ctx, id, "Chức năng dành cho quản trị viên!")
					return nil
				case "ADD_CTMS_TIMETABLE_SERVICE":
					us.AddCtmsTimetableService(ctx, id)
					return nil
				case "REMOVE_CTMS_TIMETABLE_SERVICE":
					us.RemoveCtmsTimetableService(ctx, id)
					return nil
				default:
					return nil
				}
			} else {
				log.Info().Msgf("Start chat script: %v", msg.Text)
				us.ChatScript(ctx, id, msg.Text)
			}

		}
	}
	return nil
}
