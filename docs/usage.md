# Usage Guide

This guide covers common tasks: sending notifications, scheduling jobs, and adding new notification channels.

## Sending Instant Notifications

Use the `/notify` endpoint to send immediate notifications.

### Basic Example

```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "slack",
    "recipient": "general",
    "message": {
      "title": "Alert",
      "text": "Something happened"
    }
  }'
```

### Slack Notification

```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "slack",
    "recipient": "general",
    "message": {
      "text": "Deployment completed successfully"
    }
  }'
```

### Email Notification

```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "smtp",
    "recipient": "user@example.com",
    "message": {
      "title": "Server Alert",
      "text": "CPU usage is at 85%"
    }
  }'
```

### Discord Notification

```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "discord",
    "recipient": "announcements",
    "message": {
      "text": "New release v1.2.0 is available"
    }
  }'
```

## Scheduling Recurring Notifications

Use the `/jobs` endpoint to create scheduled notifications.

### Create a Scheduled Job

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Daily Report",
    "notification_type": "smtp",
    "recipient": "team@example.com",
    "message": {
      "title": "Daily Summary",
      "text": "Your daily metrics report"
    },
    "schedule_expression": "0 9 * * *"
  }'
```

Response:
```json
{
  "id": 1,
  "name": "Daily Report",
  "notification_type": "smtp",
  "recipient": "team@example.com",
  "message": {
    "title": "Daily Summary",
    "text": "Your daily metrics report"
  },
  "schedule_expression": "0 9 * * *"
}
```

### List All Jobs

```bash
curl http://localhost:8080/jobs
```

### Cron Expression Examples

- `0 9 * * *` - Every day at 9:00 AM
- `*/15 * * * *` - Every 15 minutes
- `0 0 * * 0` - Every Sunday at midnight
- `0 9 * * 1-5` - Weekdays at 9:00 AM
- `0 */6 * * *` - Every 6 hours
- `30 8 1 * *` - First day of month at 8:30 AM

Format: `minute hour day month weekday`

## Adding New Notification Channels

The system uses a plugin architecture. Here's how to add a new channel.

### 1. Create the Plugin

Create `plugins/yourplugin/yourplugin.go`:

```go
package main

import (
    "dynamic-notification-system/config"
    "errors"
)

type YourPluginNotifier struct {
    apiKey string
}

func (n *YourPluginNotifier) Name() string {
    return "Your Plugin"
}

func (n *YourPluginNotifier) Type() string {
    return "yourplugin"
}

func (n *YourPluginNotifier) Notify(message *config.Message) error {
    // Implement your notification logic here
    return nil
}

func New(cfg map[string]interface{}) (config.Notifier, error) {
    apiKey, ok := cfg["api_key"].(string)
    if !ok || apiKey == "" {
        return nil, errors.New("missing api_key")
    }
    return &YourPluginNotifier{apiKey: apiKey}, nil
}
```

### 2. Compile the Plugin

```bash
go build -buildmode=plugin -o build/plugins/yourplugin.so plugins/yourplugin/yourplugin.go
```

### 3. Update Configuration

Add to `config.yaml`:

```yaml
channels:
  yourplugin:
    enabled: true
    api_key: "your-api-key"
```

### 4. Restart the Application

```bash
./build/notification-system
```

The plugin will be loaded automatically.

## Integration Examples

### From a Shell Script

```bash
#!/bin/bash

# Send notification on script completion
notify() {
    curl -s -X POST http://localhost:8080/notify \
      -H "Content-Type: application/json" \
      -d "{
        \"notification_type\": \"slack\",
        \"recipient\": \"ops\",
        \"message\": {
          \"text\": \"$1\"
        }
      }"
}

# Your script logic
echo "Running backup..."
# ... backup commands ...

notify "Backup completed successfully"
```

### From Python

```python
import requests

def send_notification(message):
    url = "http://localhost:8080/notify"
    payload = {
        "notification_type": "slack",
        "recipient": "general",
        "message": {
            "title": "Python Script",
            "text": message
        }
    }
    response = requests.post(url, json=payload)
    return response.status_code == 201

send_notification("Data processing complete")
```

### From Node.js

```javascript
const axios = require('axios');

async function sendNotification(message) {
    try {
        await axios.post('http://localhost:8080/notify', {
            notification_type: 'slack',
            recipient: 'general',
            message: {
                title: 'Node.js App',
                text: message
            }
        });
        console.log('Notification sent');
    } catch (error) {
        console.error('Failed to send notification:', error);
    }
}

sendNotification('Application started');
```

## Advanced Usage

### Editing Jobs

Jobs are stored in the database. You can modify them directly:

```sql
UPDATE scheduled_jobs
SET schedule_expression = '0 10 * * *'
WHERE id = 1;
```

Restart the application to reload jobs.

### Deleting Jobs

```sql
DELETE FROM scheduled_jobs WHERE id = 1;
```

Restart the application to apply changes.

Note: A DELETE endpoint will be added in a future version.

### Monitoring Job Execution

Check the `last_run` column:

```sql
SELECT id, name, last_run FROM scheduled_jobs;
```

### Debugging

Check application logs for errors:
```bash
./build/notification-system 2>&1 | tee app.log
```

Enable verbose logging by modifying the code (structured logging will be added later).

## Troubleshooting

### Notification Not Sent

Check:
- Plugin is enabled in config.yaml
- Plugin file exists in build/plugins/
- Credentials are correct
- Network connectivity to the service
- Application logs for errors

### Job Not Executing

Check:
- Scheduler is enabled: `scheduler: true`
- Job exists in database
- Cron expression is valid
- Application is running
- Check `last_run` timestamp

### Invalid Cron Expression

Use a cron validator like [crontab.guru](https://crontab.guru/) to verify your expression.

## Best Practices

- Use descriptive job names
- Test notifications before scheduling
- Monitor job execution regularly
- Keep credentials secure (use environment variables in production)
- Set up alerts for failed notifications
- Document your notification workflows
