package main

import (
	"dynamic-notification-system/config"
	"errors"
	"fmt"
)

// PushNotifier struct for push notifications
type PushNotifier struct {
	apiKey string
	device string
}

// Name returns the name of the notifier
func (p *PushNotifier) Name() string {
	return "Push Notification"
}

// Type returns the type of the notifier
func (p *PushNotifier) Type() string {
	return "push"
}

// Notify sends a push notification
// Note: Implementation incomplete - requires push service provider integration (e.g., Firebase, OneSignal)
func (p *PushNotifier) Notify(message *config.Message) error {
	return fmt.Errorf("push notification not yet implemented")
}

// New creates a new PushNotifier instance
func New(config map[string]interface{}) (config.Notifier, error) {
	apiKey, ok := config["api_key"].(string)
	device, ok2 := config["device"].(string)

	if !(ok && ok2) {
		return nil, errors.New("missing or invalid Push Notification configuration")
	}

	return &PushNotifier{
		apiKey: apiKey,
		device: device,
	}, nil
}
