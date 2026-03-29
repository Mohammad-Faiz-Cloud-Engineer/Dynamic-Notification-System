package plugins

import (
	"dynamic-notification-system/config"
	"fmt"
	"plugin"
)

func LoadPlugins(channelConfigs map[string]config.ChannelConfig) ([]config.Notifier, error) {
	var notifiers []config.Notifier

	for name, channelConfig := range channelConfigs {
		if !channelConfig.Enabled {
			continue
		}

		// Dynamically load the plugin
		plug, err := plugin.Open(fmt.Sprintf("plugins/%s.so", name))
		if err != nil {
			return nil, fmt.Errorf("error loading plugin %s: %v", name, err)
		}

		// Lookup the `New` symbol (constructor)
		sym, err := plug.Lookup("New")
		if err != nil {
			return nil, fmt.Errorf("error looking up 'New' symbol in plugin %s: %v", name, err)
		}

		// Assert the symbol's type
		constructor, ok := sym.(func(map[string]interface{}) (config.Notifier, error))
		if !ok {
			return nil, fmt.Errorf("invalid plugin constructor for %s", name)
		}

		// Convert ChannelConfig to map[string]interface{}
		configMap := map[string]interface{}{
			"enabled":      channelConfig.Enabled,
			"webhook_url":  channelConfig.WebhookURL,
			"url":          channelConfig.URL,
			"api_key":      channelConfig.ApiKey,
			"host":         channelConfig.Host,
			"port":         channelConfig.Port,
			"username":     channelConfig.Username,
			"password":     channelConfig.Password,
			"to":           channelConfig.To,
			"device":       channelConfig.Device,
			"provider_api": channelConfig.ProviderAPI,
			"phone_number": channelConfig.PhoneNumber,
			"api_url":      channelConfig.ApiURL,
			"topic":        channelConfig.Topic,
			"server":       channelConfig.Server,
		}

		// Create the notifier instance
		notifier, err := constructor(configMap)
		if err != nil {
			return nil, fmt.Errorf("error creating notifier for %s: %v", name, err)
		}

		notifiers = append(notifiers, notifier)
	}

	return notifiers, nil
}
