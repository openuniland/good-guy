package dtos

type UpdateCookieRequest struct {
	Cookie string `json:"cookie" binding:"required"`
}
