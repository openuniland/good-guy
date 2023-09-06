package middlewares

import "github.com/openuniland/good-guy/configs"

type MiddlewareManager struct {
	cfg *configs.Configs
}

func NewMiddlewareManager(cfg *configs.Configs) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg}
}
