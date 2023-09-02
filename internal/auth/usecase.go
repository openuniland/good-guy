package auth

import (
	"context"

	"github.com/openuniland/good-guy/internal/models"
)

type UseCase interface {
	Login(ctx context.Context, loginRequest *models.LoginRequest, id string) error
}
