package types

type LoginHouRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	SubscribedID string `json:"subscribed_id" validate:"required"`
}

type LoginHouResponse struct {
	Username  string `json:"username"`
	SessionId string `json:"session_id"`
	AspxAuth  string `json:"aspx_auth"`
}

type LogoutHouRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}
