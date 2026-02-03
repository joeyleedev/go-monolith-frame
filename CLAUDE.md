# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Build
go build -o bin/server cmd/server/main.go

# Run with specific environment (dev or prod)
go run cmd/server/main.go -env dev
go run cmd/server/main.go -env prod

# Format
go fmt ./...

# Vet
go vet ./...

# Download dependencies
go mod download
```

VS Code debug configuration is available at `.vscode/launch.json` - run "Launch Package" to start in dev mode.

## Architecture Overview

This is a clean architecture Go monolith using Gin + GORM. The codebase follows a layered design pattern:

- **Handler Layer** (`internal/handler/`): HTTP request handlers, handles request/response binding
- **Service Layer** (`internal/service/`): Business logic, uses repository and cache (cache-aside pattern)
- **Repository Layer** (`internal/repository/`): Data access using GORM
- **Model Layer** (`internal/model/`): Entities (database), Requests (input DTOs), Responses (output DTOs)
- **Middleware** (`internal/middleware/`): Recovery, CORS, Logger, Auth (JWT)
- **Infrastructure** (`pkg/`): Reusable packages - cache, mysql, logger, response, utils (JWT, password)

## Application Startup Sequence

`cmd/server/main.go` initializes components in order:

1. Parse CLI flags (`-env` for dev/prod)
2. Load YAML config from `configs/config.{env}.yaml`
3. Initialize logger (zap + lumberjack for rotation)
4. Initialize MySQL database with GORM (auto-migrate runs for all entities)
5. Initialize cache (memory or redis based on config)
6. Wire dependencies: Repository -> Service -> Handler
7. Setup Gin router with middleware stack
8. Register routes (public routes, then protected routes with Auth middleware)
9. Start HTTP server with graceful shutdown on SIGINT/SIGTERM

## Configuration Structure

Config files in `configs/`:

- `config.dev.yaml` - Development (memory cache, debug mode, console logs, 24h JWT expiry)
- `config.prod.yaml` - Production (redis cache, release mode, file-only logs, 7d JWT expiry)

Key sections:
- `server`: addr, mode (debug/release), timeouts
- `database`: dsn, pool settings, log_level, slow_threshold
- `cache`: type (memory/redis), backend-specific settings
- `log`: level, filename, rotation settings, console output
- `jwt`: secret, expire_time

## Routing Pattern

Base URL: `/api/v1`

Middleware order (applied globally): Recovery -> CORS -> Logger

Public routes (no auth):
- `POST /api/v1/users/login` - Returns JWT token
- `POST /api/v1/users` - Create user

Protected routes (JWT auth required via `middleware.Auth()`):
- `GET /api/v1/users` - List with pagination and search
- `GET /api/v1/users/:id` - Get by ID
- `PUT /api/v1/users/:id` - Update
- `DELETE /api/v1/users/:id` - Delete
- `PATCH /api/v1/users/:id/password` - Change password

Health check: `GET /health`

## Error Code Convention

Error codes use prefixes by module:
- `0xxxx` - System/common errors (defined in `pkg/errcode/`)
- `2xxxx` - User module errors (defined in `internal/errcode/`)

Example: `20001` - User not found, `20005` - Invalid token

## Adding a New Feature

1. Define entity in `internal/model/entity/` with GORM tags
2. Define request/response DTOs in `internal/model/request/` and `internal/model/response/`
3. Create repository in `internal/repository/` (implements data access)
4. Create service in `internal/service/` (business logic, receives repo and cache)
5. Create handler in `internal/handler/` (HTTP layer, receives service)
6. Wire dependencies in `cmd/server/main.go` (repo -> service -> handler)
7. Register routes: add to public or protected route group

## Database Notes

- GORM auto-migration runs on startup for all entities in `internal/model/entity/`
- Entity uses soft deletes via `gorm.DeletedAt`
- Custom Zap-based GORM logger at `pkg/mysql/logger.go` supports slow query detection
- Connection pooling configured in yaml: max_open_conns, max_idle_conns, conn_max_lifetime, conn_max_idle_time

## Cache Pattern

Cache-aside pattern used in services:
- Read: Check cache first, if miss fetch from DB and populate cache
- Write: Delete/invalidate cache on create/update/delete
- Cache key format: `{entity}:{id}` (e.g., `user:1`)
- TTL configured per use case (UserService uses 10 minutes)
