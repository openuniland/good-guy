package usecase

import (
	"context"
	"errors"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/constants"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/external/types"
	"github.com/openuniland/good-guy/internal/common"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/openuniland/good-guy/pkg/utils"
	"github.com/rs/zerolog/log"
)

type CommonUS struct {
	cfg        *configs.Configs
	facebookUS facebook.UseCase
	ctmsUC     ctms.UseCase
	userUC     users.UseCase
}

func NewCommonUseCase(cfg *configs.Configs, facebookUS facebook.UseCase, ctmsUC ctms.UseCase, userUC users.UseCase) common.UseCase {
	return &CommonUS{cfg: cfg, facebookUS: facebookUS, ctmsUC: ctmsUC, userUC: userUC}
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

func (us *CommonUS) HandleFacebookWebhook(ctx context.Context, data *types.FacebookWebhookRequest) error {
	messaging := data.Entry[0].Messaging

	for _, message := range messaging {
		sender := message.Sender
		postback := message.Postback
		msg := message.Message

		id := sender.ID

		if postback.Payload != "" {
			switch postback.Payload {
			case "GET_STARTED":
				us.facebookUS.SendTextMessage(ctx, id, "Chào mừng bạn đến với Fithou BOT. Chúc bạn có một trải nghiệm zui zẻ :D. /help để biết thêm chi tiết!")
				return nil
			case "HELP":
				us.facebookUS.SendTextMessage(ctx, id, utils.HelpScript())
				return nil
			case "CTMS_SERVICE":
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReply{{
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
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReply{{
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
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReply{{
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
				us.facebookUS.SendQuickReplies(ctx, id, "Chọn một câu trả lời:", &[]types.QuickReply{{
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
		} else if msg.QuickReply.Payload != "" {
			switch msg.QuickReply.Payload {
			case "ADD_CTMS_ACCOUNT":
				// TODO: send btn login
				return nil
			case "REMOVE_CTMS_ACCOUNT":
				// TODO: send btn remove
				return nil
			case "ADD_FITHOU_CRAWL_SERVICE":
				// TODO: subscribe fithou notification
				return nil
			case "REMOVE_FITHOU_CRAWL_SERVICE":
				// TODO: unsubscribe fithou notification
				return nil
			case "ADD_CTMS_CREDITS_SERVICE":
				//
				return nil
			case "REMOVE_CTMS_CREDITS_SERVICE":
				return nil
			case "ADD_CTMS_TIMETABLE_SERVICE":
				return nil
			case "REMOVE_CTMS_TIMETABLE_SERVICE":
				return nil
			default:
				return nil
			}
		} else if msg.Text != "" {
			utils.ChatScript(id, msg.Text)

		}
	}
	return nil
}
