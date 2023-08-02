package ctms

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	Login(ctx context.Context, user *types.LoginRequest) (*types.LoginResponse, error)
	Logout(ctx context.Context, cookie string) error
	GetDailySchedule(ctx context.Context, cookie string) ([]*types.DailySchedule, error)
}
