# Dynamic Notification System

A flexible, plugin-based notification system that lets you send alerts across multiple platforms from a single API. Built with Go for performance and reliability.

[![Build and Test](https://github.com/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System/actions/workflows/build.yml/badge.svg)](https://github.com/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System/actions/workflows/build.yml) ![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square) ![Go Version](https://img.shields.io/badge/go-1.23-blue?logo=go&style=flat-square) ![GitHub release](https://img.shields.io/github/v/release/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System?style=flat-square)

## What is this?

This is a notification service that solves a common problem: sending alerts to different platforms without writing custom code for each one. Instead of maintaining separate integrations for Slack, Discord, email, SMS, and other services, you configure them once and send notifications through a simple REST API.

The system uses a plugin architecture, so adding support for a new platform is straightforward. It also includes a scheduler for recurring notifications and stores everything in MySQL.

## Key Features

**Plugin Architecture** - Each notification platform is a separate plugin. Enable only what you need, and the system loads them dynamically at startup.

**Scheduled Notifications** - Set up recurring alerts using cron expressions. The scheduler runs in the background and executes jobs automatically.

**REST API** - Send instant notifications or manage scheduled jobs through HTTP endpoints. Simple JSON payloads, no complex setup.

**Configuration-Based** - All settings live in a YAML file. No hardcoded values, easy to deploy across environments.

**Production Ready** - Includes connection pooling, HTTP timeouts, parameterized SQL queries, and proper error handling.

## Supported Platforms

Currently implemented:
- Slack (webhooks)
- Discord (webhooks)
- Microsoft Teams (webhooks)
- Email (SMTP)
- Rocket.Chat (webhooks)
- Generic Webhooks
- Ntfy (push notifications)

In progress:
- Telegram (API integration needed)
- SMS (provider integration needed)
- Signal (API integration needed)
- Push notifications (Firebase/OneSignal integration needed)

## Why I Built This

I needed a way to centralize notifications for multiple projects without duplicating code. Each project had its own notification logic scattered across the codebase, making it hard to maintain and extend. This system extracts that complexity into a standalone service.

The plugin design means you can add new platforms without touching the core code. The scheduler handles recurring alerts without external cron jobs. And the REST API makes it easy to integrate with any application.

## Quick Start

### Prerequisites

- Go 1.23 or higher
- MySQL 8.0 or higher
- Docker (optional, for containerized deployment)

### Installation

Clone the repository:
```bash
git clone https://github.com/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System.git
cd Dynamic-Notification-System
```

Set up your configuration:
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your credentials
```

Initialize the database:
```bash
mysql -u root -p < db/init.sql
```

Build the application:
```bash
make all
```

Run the server:
```bash
./build/notification-system
```

The server starts on port 8080 by default.

### Docker Deployment

If you prefer Docker:
```bash
docker-compose up -d
```

This starts both the application and MySQL database. The config file is mounted as a volume, so you can edit it without rebuilding.

## Configuration

Copy `config.yaml.example` to `config.yaml` and update with your settings:

```yaml
scheduler: true

database:
  user: "your_user"
  password: "your_password"
  host: "localhost"
  port: 3306
  name: "notifications"

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

Only enabled channels are loaded. The system validates configuration at startup and fails fast if something is wrong.

## Usage

### Send an Instant Notification

```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "slack",
    "recipient": "general",
    "message": {
      "title": "Server Alert",
      "message": "CPU usage above 80%"
    }
  }'
```

### Schedule a Recurring Notification

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Daily Report",
    "notification_type": "smtp",
    "recipient": "team@example.com",
    "message": {
      "title": "Daily Summary",
      "message": "Your daily metrics report"
    },
    "schedule_expression": "0 9 * * *"
  }'
```

The schedule expression uses standard cron format:
- `0 9 * * *` - Every day at 9 AM
- `*/15 * * * *` - Every 15 minutes
- `0 0 * * 0` - Every Sunday at midnight

### List Scheduled Jobs

```bash
curl http://localhost:8080/jobs
```

### Get Job Schema

```bash
curl http://localhost:8080/schema/job
```

Returns a JSON schema describing the job structure, useful for validation.

## Project Structure

```
.
├── config/          # Configuration loading and types
├── db/              # Database initialization scripts
├── notifier/        # Instant notification handler
├── plugins/         # Notification platform plugins
│   ├── slack/
│   ├── discord/
│   ├── smtp/
│   └── ...
├── scheduler/       # Job scheduling and cron management
├── main.go          # Application entry point
└── config.yaml      # Configuration file (not in repo)
```

Each plugin is compiled as a shared object (.so file) and loaded dynamically based on your configuration.

## Security Considerations

- Never commit `config.yaml` to version control (it's in .gitignore)
- Use `config.yaml.example` as a template
- Restrict file permissions: `chmod 600 config.yaml`
- Run behind a reverse proxy with TLS in production
- Consider adding API authentication for public deployments
- See [SECURITY.md](SECURITY.md) for detailed security guidelines

## Known Limitations

- No authentication on API endpoints (add reverse proxy auth or implement middleware)
- No rate limiting (vulnerable to abuse without external protection)
- Some plugins are incomplete (Telegram, SMS, Signal, Push)
- No test coverage yet
- No graceful shutdown handling

These are documented in [AUDIT_REPORT.md](AUDIT_REPORT.md) with recommendations.

## Contributing

Contributions are welcome. Here's how you can help:

- Implement the incomplete plugins (Telegram, SMS, Signal, Push)
- Add new notification platforms
- Improve error handling and logging
- Write tests (currently none exist)
- Enhance documentation
- Report bugs or suggest features

See [docs/contributing.md](docs/contributing.md) for detailed guidelines.

## Development

Build the main application:
```bash
make build
```

Build plugins:
```bash
make build-plugins
```

Clean build artifacts:
```bash
make clean
```

The Makefile includes targets for common tasks. Check `make help` for all options.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

Created and maintained by [Mohammad Faiz](https://github.com/Mohammad-Faiz-Cloud-Engineer).

If you find this useful, consider starring the repository or contributing improvements.
