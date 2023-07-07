package types

type SendMessageRequest struct {
	Text string `json:"text" validate:"required"`
}

type SendButtonMessageRequest struct {
	ImageUrl string `json:"image_url" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Subtitle string `json:"subtitle" validate:"required"`
	Url      string `json:"url" validate:"required"`
	BtnText  string `json:"btn_text" validate:"required"`
}
