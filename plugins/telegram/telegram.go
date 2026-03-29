package main

import (
	"dynamic-notification-system/config"
	"errors"
	"fmt"
)

type TelegramNotifier struct {
	apiKey string
}

func (t *TelegramNotifier) Name() string {
	return "Telegram"
}

func (t *TelegramNotifier) Type() string {
	return "telegram"
}

// Notify sends a message to Telegram
// Note: Implementation incomplete - requires Telegram Bot API integration
func (t *TelegramNotifier) Notify(message *config.Message) error {
	return fmt.Errorf("telegram notification not yet implemented")
}

func New(config map[string]interface{}) (config.Notifier, error) {
	apiKey, ok := config["api_key"].(string)
	if !ok || apiKey == "" {
		return nil, errors.New("missing or invalid API key")
	}
	return &TelegramNotifier{apiKey: apiKey}, nil
}
