package usecase

import (
	"context"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/auth"
	"github.com/openuniland/good-guy/internal/models"
)

type AuthUS struct {
	cfg *configs.Configs
}

func NewAuthUseCase(cfg *configs.Configs) auth.UseCase {
	return &AuthUS{cfg: cfg}
}

func (a *AuthUS) Login(ctx context.Context, loginRequest *models.LoginRequest, id string) error {

	return nil
}
