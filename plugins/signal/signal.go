package main

import (
	"dynamic-notification-system/config"
	"errors"
	"fmt"
)

// SignalNotifier struct for Signal messaging
type SignalNotifier struct {
	phoneNumber string
	apiURL      string
}

// Name returns the name of the notifier
func (s *SignalNotifier) Name() string {
	return "Signal"
}

// Type returns the type of the notifier
func (s *SignalNotifier) Type() string {
	return "signal"
}

// Notify sends a message via Signal
// Note: Implementation incomplete - requires Signal API integration
func (s *SignalNotifier) Notify(message *config.Message) error {
	return fmt.Errorf("signal notification not yet implemented")
}

// New creates a new SignalNotifier instance
func New(config map[string]interface{}) (config.Notifier, error) {
	phoneNumber, ok := config["phone_number"].(string)
	apiURL, ok2 := config["api_url"].(string)

	if !(ok && ok2) {
		return nil, errors.New("missing or invalid Signal configuration")
	}

	return &SignalNotifier{
		phoneNumber: phoneNumber,
		apiURL:      apiURL,
	}, nil
}
