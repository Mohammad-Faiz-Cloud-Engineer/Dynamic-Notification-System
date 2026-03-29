# Contributing

Thanks for your interest in contributing to the Dynamic Notification System. This document explains how to contribute effectively.

## Getting Started

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
```bash
git clone https://github.com/your-username/Dynamic-Notification-System.git
cd Dynamic-Notification-System
```

### Set Up Development Environment

1. Install Go 1.23 or higher
2. Install MySQL 8.0 or higher
3. Install dependencies:
```bash
go mod download
```

### Create a Branch

Use descriptive branch names:
```bash
git checkout -b feature/add-telegram-plugin
git checkout -b fix/scheduler-memory-leak
git checkout -b docs/improve-api-examples
```

## Making Changes

### Code Style

Follow standard Go conventions:
- Run `gofmt -s -w .` before committing
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Add comments for exported functions and complex logic
- Keep functions focused and under 50 lines when possible

### Testing

Currently, the project lacks tests. If you're adding tests (which would be great):
- Place test files next to the code they test
- Use table-driven tests for multiple scenarios
- Mock external dependencies (database, HTTP calls)

### Commit Messages

Write clear commit messages:
```bash
git commit -m "Add Telegram notification plugin"
git commit -m "Fix race condition in scheduler"
git commit -m "Update API documentation with examples"
```

For larger changes, include more detail:
```
Add Telegram notification plugin

- Implement Notifier interface
- Add configuration parsing
- Include error handling for API failures
- Update documentation
```

## Submitting Changes

### Push Your Branch

```bash
git push origin feature/your-feature-name
```

### Open a Pull Request

1. Go to the original repository on GitHub
2. Click "New Pull Request"
3. Select your branch
4. Write a clear description of your changes:
   - What problem does it solve?
   - How did you test it?
   - Are there any breaking changes?

### Code Review

Expect feedback on your pull request. Common review points:
- Code style and formatting
- Error handling
- Security considerations
- Documentation updates
- Test coverage (when tests exist)

Be responsive to feedback and make requested changes promptly.

## What to Contribute

### High Priority

- Complete incomplete plugins (Telegram, SMS, Signal, Push)
- Add unit and integration tests
- Implement API authentication
- Add request body size limits
- Implement graceful shutdown

### Medium Priority

- Add health check endpoints
- Implement rate limiting
- Add job deletion endpoint
- Improve error messages
- Add structured logging

### Low Priority

- Add Prometheus metrics
- Implement CORS configuration
- Add OpenAPI/Swagger documentation
- Improve plugin error handling
- Add more notification platforms

## Reporting Issues

Found a bug or have a feature request?

1. Check if an issue already exists
2. If not, open a new issue with:
   - Clear title
   - Detailed description
   - Steps to reproduce (for bugs)
   - Expected vs actual behavior
   - Your environment (OS, Go version, etc.)

## Code of Conduct

Be respectful and professional in all interactions. We're all here to build something useful.

## Questions?

Open an issue with the "question" label or start a discussion on GitHub.
