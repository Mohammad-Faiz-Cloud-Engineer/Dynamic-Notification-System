# Dynamic Notification System Documentation

## Overview

The Dynamic Notification System is a Go-based service for managing notifications across multiple platforms. It provides a REST API for sending instant notifications and scheduling recurring alerts using cron expressions.

Key capabilities:
- Send notifications to Slack, Discord, Teams, Email, and other platforms
- Schedule recurring notifications with cron expressions
- Plugin-based architecture for easy extensibility
- MySQL database for job persistence
- RESTful API for integration

## Architecture

### Core Components

**Configuration Loader**
- Reads settings from `config.yaml`
- Validates database and channel configurations
- Fails fast on startup if configuration is invalid

**Database Layer**
- Stores scheduled job metadata
- Tracks execution history
- Uses parameterized queries to prevent SQL injection

**Plugin System**
- Each notification platform is a separate plugin
- Plugins are compiled as shared objects (.so files)
- Loaded dynamically at startup based on configuration
- Implements a common `Notifier` interface

**Scheduler**
- Built on `robfig/cron/v3` library
- Executes jobs based on cron expressions
- Runs in a background goroutine
- Updates last_run timestamp after each execution

**REST API**
- Built with Gorilla Mux router
- Handles instant notifications and job management
- Returns JSON responses
- Includes proper error handling

## Setup and Installation

### Prerequisites
- Go 1.23 or higher
- MySQL 8.0 or higher
- Linux environment (for plugin support)

### Installation Steps

1. Clone the repository:
```bash
git clone https://github.com/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System.git
cd Dynamic-Notification-System
```

2. Build the application:
```bash
make all
```

3. Set up configuration:
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your credentials
```

4. Initialize the database:
```bash
mysql -u root -p < db/init.sql
```

5. Run the application:
```bash
./build/notification-system
```

The server starts on port 8080.

## Configuration Example

```yaml
scheduler: true

database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  name: notifications

channels:
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
  
  smtp:
    enabled: true
    host: "smtp.gmail.com"
    port: "587"
    username: "your-email@gmail.com"
    password: "your-app-password"
    to: "recipient@example.com"
```

## Documentation Index

1. [Getting Started](getting_started.md) - Installation and first steps
2. [Usage Guide](usage.md) - How to use the API and add plugins
3. [Developer Guide](developer_guide.md) - Architecture and plugin development
4. [Contributing](contributing.md) - How to contribute to the project

Technical documentation:
- [Main Module](technical_docs/main.md)
- [Scheduler Module](technical_docs/scheduler.md)
- [Notifier Module](technical_docs/notifier.md)
- [Config Module](technical_docs/config.md)

## API Reference

### POST /notify
Send an instant notification.

Request:
```json
{
  "notification_type": "slack",
  "recipient": "general",
  "message": {
    "title": "Server Alert",
    "text": "CPU usage above 80%"
  }
}
```

Response: 201 Created

### POST /jobs
Create a scheduled notification job.

Request:
```json
{
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

### GET /jobs
Retrieve all scheduled jobs.

Response:
```json
[
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
]
```

### GET /schema/job
Get the JSON schema for job objects. Useful for validation.

## Dependencies

- `github.com/gorilla/mux` - HTTP routing
- `github.com/robfig/cron/v3` - Cron-based scheduling
- `github.com/go-sql-driver/mysql` - MySQL driver
- `gopkg.in/yaml.v3` - YAML configuration parsing

## Roadmap

Planned improvements:
- Complete Telegram, SMS, Signal, and Push notification plugins
- Add API authentication (API keys or JWT)
- Implement rate limiting
- Add health check endpoints
- Write comprehensive test suite
- Add Prometheus metrics
- Implement graceful shutdown

## License

MIT License - see [LICENSE](https://github.com/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System/blob/main/LICENSE)

## Author

Created by [Mohammad Faiz](https://github.com/Mohammad-Faiz-Cloud-Engineer)
