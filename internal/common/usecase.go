package common

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	VerifyFacebookWebhook(ctx context.Context, token, challenge string) (string, error)
	HandleFacebookWebhook(ctx context.Context, data *types.FacebookWebhookRequest) error
}
