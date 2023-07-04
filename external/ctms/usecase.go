package ctms

import (
	"context"

	"github.com/openuniland/good-guy/external/models"
)

type UseCase interface {
	Login(ctx context.Context, user *models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, cookie string) error
}
