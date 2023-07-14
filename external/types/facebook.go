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

type FacebookWebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Postback *struct {
				Title   string `json:"title"`
				Payload string `json:"payload"`
			} `json:"postback"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   *struct {
				Mid        string `json:"mid"`
				Text       string `json:"text"`
				QuickReply struct {
					Payload string `json:"payload"`
				} `json:"quick_reply"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}
