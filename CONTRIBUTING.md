# Contributing to Batch-GPT

Thank you for your interest in contributing to Batch-GPT! This document provides guidelines and best practices for contributions.

## Development Environment

1. Fork and clone the repository
2. Install prerequisites:
   - Go 1.23.0 or later
   - Docker and Docker Compose
   - Python 3.x (for test client)
3. Set up MongoDB:
   ```bash
   cd local/mongo
   docker-compose up -d
   ```
4. Set required environment variables as described in README.md

## Project Structure

- `/server`: Main server code
  - `/db`: Database interactions
  - `/handlers`: HTTP request handlers
  - `/logger`: Custom logging setup
  - `/models`: Data models
  - `/services`: Business logic and OpenAI interactions
    - `/batch`: Batch processing logic
    - `/cache`: Caching logic
    - `/client`: OpenAI client wrapper
    - `/config`: Configuration management
    - `/utils`: Common utilities
- `/cmd`: Command-line tools
  - `/monitor`: Terminal-based monitoring tool
    - `/ui`: UI components and styling
- `/local/mongo`: Local MongoDB setup
- `/test-python-client`: Test client

## Adding New Features

1. Choose the appropriate package:
   - `/server/handlers` for new API endpoints
   - `/server/db` for database interactions
   - `/server/services/*` for business logic
   - `/cmd/monitor/ui` for monitoring UI components

2. Implementation:
   - Add types/interfaces in `types.go` if needed
   - Implement core functionality
   - Update services/orchestrators as needed
   - Add environment variables to config package

3. Integration:
   - Wire new handlers in `server/main.go`
   - Update UI model in `ui/model.go` for UI changes
   - Create necessary database indexes

4. Testing:
   - Test in all serving modes (sync/async/cache)
   - Update test Python client for API changes

## Code Style

- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and concise
- Write idiomatic Go code

## Pull Request Process

1. Create a feature/fix branch (`feature/description` or `fix/description`)
2. Make focused, incremental changes
3. Test changes in all serving modes
4. Update documentation as needed
5. Submit PR with clear description of changes
6. Respond to review comments

## Issue Reporting

When creating an issue, include:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

## Questions?

Feel free to open an issue for any questions about contributing.
