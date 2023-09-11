package ctms

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	LoginCtms(ctx context.Context, user *types.LoginCtmsRequest) (*types.LoginCtmsResponse, error)
	LogoutCtms(ctx context.Context, cookie string) error
	GetDailySchedule(ctx context.Context, cookie string) ([]*types.DailySchedule, error)
	GetExamSchedule(ctx context.Context, cookie string) ([]types.ExamSchedule, error)
	GetUpcomingExamSchedule(ctx context.Context, user *types.LoginCtmsRequest) (types.GetUpcomingExamScheduleResponse, error)
	SendChangedExamScheduleAndNewExamScheduleToUser(ctx context.Context, user *types.LoginCtmsRequest, id string) error
	ForceLogout(ctx context.Context, subscribed_id string) error
}
