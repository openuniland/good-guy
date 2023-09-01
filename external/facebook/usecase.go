package facebook

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	SendMessage(ctx context.Context, id string, message interface{}) error
	SendQuickReplies(ctx context.Context, id string, message string, quickReplies *[]types.QuickReply) error
	SendButtonMessage(ctx context.Context, id string, input *types.SendButtonMessageRequest) error
	SendTextMessage(ctx context.Context, id string, text string) error
	VerifyFacebookWebhook(ctx context.Context, token, challenge string) (string, error)
	HandleFacebookWebhook(ctx context.Context, data *types.FacebookWebhookRequest) error
}
