# Developer Guide

This guide covers the architecture, plugin development, and deployment of the Dynamic Notification System.

## Architecture Overview

### Main Application

The entry point (`main.go`) handles:
- Loading configuration from `config.yaml`
- Initializing database connection with pooling
- Loading enabled notification plugins
- Starting the scheduler (if enabled)
- Setting up HTTP server with timeouts
- Routing API requests

### Plugin System

Plugins are Go shared objects (.so files) that implement the `Notifier` interface:

```go
type Notifier interface {
    Name() string                    // Human-readable name
    Type() string                    // Type identifier (matches config)
    Notify(message *Message) error   // Send notification
}
```

Each plugin:
- Lives in its own directory under `plugins/`
- Compiles to a `.so` file
- Exports a `New` constructor function
- Receives configuration as a map

The plugin loader:
- Reads enabled channels from config
- Opens corresponding `.so` files
- Looks up the `New` symbol
- Calls the constructor with config
- Stores the notifier instance

### Scheduler

The scheduler uses `robfig/cron/v3` to execute jobs:

1. On startup, loads all jobs from the database
2. Adds each job to the cron instance
3. Runs in a background goroutine
4. When a job triggers:
   - Finds the matching notifier by type
   - Calls `Notify()` with the message
   - Updates `last_run` timestamp in database
   - Logs any errors

### Database

MySQL stores scheduled jobs with this schema:

```sql
CREATE TABLE scheduled_jobs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    notification_type VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    message TEXT,
    schedule_expression VARCHAR(255) NOT NULL,
    last_run DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

The `message` column stores JSON. The `config` package provides `Scan()` and `Value()` methods for automatic marshaling.

Connection pooling is configured:
- MaxOpenConns: 25
- MaxIdleConns: 5
- ConnMaxLifetime: 5 minutes

## Plugin Development

### Creating a New Plugin

1. Create a directory: `plugins/yourplugin/`

2. Create `yourplugin.go`:

```go
package main

import (
    "dynamic-notification-system/config"
    "errors"
    "fmt"
)

type YourPluginNotifier struct {
    apiKey string
    // other config fields
}

func (n *YourPluginNotifier) Name() string {
    return "Your Plugin"
}

func (n *YourPluginNotifier) Type() string {
    return "yourplugin"  // Must match config key
}

func (n *YourPluginNotifier) Notify(message *config.Message) error {
    // Implement notification logic
    // Use message.Title, message.Text, etc.
    return nil
}

// Constructor - required for plugin loading
func New(cfg map[string]interface{}) (config.Notifier, error) {
    apiKey, ok := cfg["api_key"].(string)
    if !ok || apiKey == "" {
        return nil, errors.New("missing api_key")
    }
    
    return &YourPluginNotifier{
        apiKey: apiKey,
    }, nil
}
```

3. Compile the plugin:

```bash
go build -buildmode=plugin -o build/plugins/yourplugin.so plugins/yourplugin/yourplugin.go
```

4. Add to `config.yaml`:

```yaml
channels:
  yourplugin:
    enabled: true
    api_key: "your-api-key"
```

5. Update the Makefile to include your plugin in `build-plugins` target.

### Plugin Best Practices

- Validate all configuration in the `New` constructor
- Return descriptive errors
- Use HTTP clients with timeouts (30 seconds recommended)
- Don't log sensitive data (API keys, tokens)
- Handle rate limiting gracefully
- Implement retries for transient failures
- Keep the `Notify` method focused and simple

### Testing Plugins

Since plugins are separate binaries, test them by:
1. Building the plugin
2. Running the main application
3. Sending a test notification via the API
4. Checking logs for errors

Eventually, add unit tests that mock the plugin interface.

## Database Schema

### Scheduled Jobs Table

```sql
CREATE TABLE scheduled_jobs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    notification_type VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    message TEXT,
    schedule_expression VARCHAR(255) NOT NULL,
    last_run DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Fields:
- `id`: Auto-incrementing primary key
- `name`: Human-readable job name
- `notification_type`: Must match a plugin's `Type()`
- `recipient`: Platform-specific recipient (email, channel, phone, etc.)
- `message`: JSON-encoded message object
- `schedule_expression`: Cron expression (e.g., "0 9 * * *")
- `last_run`: Timestamp of last execution
- `created_at`: Job creation timestamp

## Deployment

### Docker

Build the image:
```bash
docker build -t dynamic-notification-system .
```

Run with docker-compose:
```bash
docker-compose up -d
```

The compose file includes:
- Application container
- MySQL container
- Volume for database persistence
- Volume mount for config.yaml

### Manual Deployment

1. Build the application:
```bash
make all
```

2. Copy files to server:
```bash
scp -r build/ config.yaml user@server:/opt/notification-system/
```

3. Set up systemd service (example):

```ini
[Unit]
Description=Dynamic Notification System
After=network.target mysql.service

[Service]
Type=simple
User=notifier
WorkingDirectory=/opt/notification-system
ExecStart=/opt/notification-system/build/notification-system
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

4. Start the service:
```bash
sudo systemctl enable notification-system
sudo systemctl start notification-system
```

### Production Considerations

- Run behind a reverse proxy (nginx, Caddy) with TLS
- Add authentication (API keys, JWT)
- Implement rate limiting at the proxy level
- Set up log rotation
- Configure monitoring and alerting
- Use a secrets manager for credentials
- Set up automated backups for the database

## API Endpoints

### POST /notify
Send an instant notification.

### POST /jobs
Create a scheduled job.

### GET /jobs
List all scheduled jobs.

### GET /schema/job
Get the JSON schema for job validation.

See the [main documentation](index.md) for detailed API examples.

## Troubleshooting

### Plugin Not Loading

Check:
- Plugin file exists in `build/plugins/`
- Plugin name matches config key
- Plugin exports a `New` function
- Configuration is valid

### Scheduler Not Running

Check:
- `scheduler: true` in config.yaml
- Database connection is working
- Jobs exist in the database
- Cron expressions are valid

### Database Connection Fails

Check:
- MySQL is running
- Credentials are correct
- Database exists
- Network connectivity
- Connection string format

## Further Reading

- [Go Plugin Package](https://pkg.go.dev/plugin)
- [Cron Expression Format](https://pkg.go.dev/github.com/robfig/cron/v3)
- [Gorilla Mux Documentation](https://github.com/gorilla/mux)
