package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/rs/zerolog/log"
)

type FacebookUS struct {
	cfg *configs.Configs
}

func NewCtmsUseCase(cfg *configs.Configs) facebook.UseCase {
	return &FacebookUS{cfg: cfg}
}

func (us *FacebookUS) SendMessage(ctx context.Context, id string, message interface{}) error {
	data := map[string]interface{}{
		"recipient": map[string]string{
			"id": id,
		},
		"message":        message,
		"messaging_type": "MESSAGE_TAG",
		"tag":            "ACCOUNT_UPDATE",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("error marshal data")
		return err
	}

	url := fmt.Sprintf("https://graph.facebook.com/v14.0/me/messages?access_token=%s", us.cfg.FBConfig.FBVerifyToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("error create new request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error send request")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Err(err).Msg("error response from Facebook API")
		return err
	}

	return nil
}
