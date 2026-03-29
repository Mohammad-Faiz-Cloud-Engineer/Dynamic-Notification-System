# Getting Started

This guide walks you through installing and running the Dynamic Notification System.

## Prerequisites

### Required

- **Go 1.23 or higher**
  - Download from [golang.org](https://golang.org/dl/)
  - Verify: `go version`

- **MySQL 8.0 or higher**
  - Install for your platform or use Docker
  - Verify: `mysql --version`

### Optional

- **Docker and Docker Compose**
  - For containerized deployment
  - Download from [docker.com](https://www.docker.com/products/docker-desktop)
  - Verify: `docker --version` and `docker-compose --version`

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Mohammad-Faiz-Cloud-Engineer/Dynamic-Notification-System.git
cd Dynamic-Notification-System
```

### 2. Set Up Configuration

Copy the example config:
```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml` with your settings:
```yaml
scheduler: true

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"
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

Enable only the channels you need. Disabled channels won't be loaded.

### 3. Initialize the Database

Create the database and load the schema:
```bash
mysql -u root -p -e "CREATE DATABASE notifications;"
mysql -u root -p notifications < db/init.sql
```

This creates the `scheduled_jobs` table and adds sample data.

### 4. Build the Application

```bash
make all
```

This compiles:
- The main application to `build/notification-system`
- All plugins to `build/plugins/*.so`

### 5. Run the Application

```bash
./build/notification-system
```

You should see:
```
Server listening on port :8080
```

The application is now running and ready to accept requests.

## Docker Deployment

If you prefer Docker:

### 1. Set Up Configuration

```bash
cp config.yaml.example config.yaml
# Edit config.yaml - use "db" as the database host
```

### 2. Start Services

```bash
docker-compose up -d
```

This starts:
- MySQL database on port 3306
- Application on port 8080

### 3. Check Status

```bash
docker-compose ps
docker-compose logs -f app
```

### 4. Stop Services

```bash
docker-compose down
```

To remove data volumes:
```bash
docker-compose down -v
```

## Verify Installation

### Check Server Health

```bash
curl http://localhost:8080/jobs
```

Should return an empty array or existing jobs.

### Send a Test Notification

If you configured Slack:
```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "slack",
    "recipient": "general",
    "message": {
      "title": "Test",
      "text": "Hello from Dynamic Notification System"
    }
  }'
```

Check your Slack channel for the message.

## Configuration Tips

### Database Connection

For local development:
```yaml
database:
  host: "localhost"
  port: 3306
```

For Docker Compose:
```yaml
database:
  host: "db"  # Service name from docker-compose.yml
  port: 3306
```

### Email Configuration

For Gmail, use an app password (not your regular password):
1. Enable 2-factor authentication
2. Generate an app password at [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)
3. Use the app password in config.yaml

### Webhook URLs

Get webhook URLs from:
- **Slack**: Workspace Settings → Apps → Incoming Webhooks
- **Discord**: Server Settings → Integrations → Webhooks
- **Teams**: Channel → Connectors → Incoming Webhook

## Troubleshooting

### "Error loading config"

Check:
- `config.yaml` exists in the current directory
- YAML syntax is valid (use a YAML validator)
- All required fields are present

### "Failed to connect to DB"

Check:
- MySQL is running
- Database exists
- Credentials are correct
- Host and port are correct

### "Error loading plugin"

Check:
- Plugin files exist in `build/plugins/`
- You ran `make build-plugins`
- Plugin name in config matches the .so filename

### Port 8080 Already in Use

Change the port in `main.go`:
```go
const serverPort = ":8081"  // Or any available port
```

Then rebuild: `make build`

## Next Steps

- [Usage Guide](usage.md) - Learn how to use the API
- [Developer Guide](developer_guide.md) - Understand the architecture
- [Contributing](contributing.md) - Help improve the project
