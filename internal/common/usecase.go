package common

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	VerifyFacebookWebhook(ctx context.Context, token, challenge string) (string, error)
	HandleFacebookWebhook(ctx context.Context, data *types.FacebookWebhookRequest) error
	SendLoginCtmsButton(ctx context.Context, id string) error
	RemoveUser(ctx context.Context, id string) error
	AddFithouCrawlService(ctx context.Context, id string) error
	RemoveFithouCrawlService(ctx context.Context, id string) error
	AddCtmsTimetableService(ctx context.Context, id string) error
	RemoveCtmsTimetableService(ctx context.Context, id string) error
	ChatScript(ctx context.Context, id string, msg string)
	GetNotificationOfExamSchedule(ctx context.Context, id string) error
	CancelGetNotificationOfExamSchedule(ctx context.Context, id string) error
}
