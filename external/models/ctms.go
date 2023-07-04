package models

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Cookie    string `json:"cookie"`
	IsSuccess bool   `json:"is_success"`
}

type LogoutRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}
