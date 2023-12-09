package hou

import (
	"context"

	"github.com/openuniland/good-guy/external/types"
)

type UseCase interface {
	LoginHou(ctx context.Context, user *types.LoginHouRequest) (*types.LoginHouResponse, error)
	LogoutHou(ctx context.Context, SessionId string) error
}
