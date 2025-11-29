# Design Documentation

## Overview

This document describes the architecture and design decisions for the Go REST API Template.

## Architecture

### Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and queries
│   ├── handler/         # HTTP request handlers
│   └── server/          # HTTP server setup
├── api/                 # API specifications (OpenAPI)
├── db/                  # Database files
│   ├── queries/         # sqlc query files
│   ├── schema.sql       # Database schema
│   └── sample-data.sql  # Seed data
└── doc/                 # Documentation
```

### Components

#### Configuration (internal/config)

- Uses **Viper** for configuration management
- Supports multiple sources with precedence: CLI flags > environment variables > defaults
- Environment variables use `APP_` prefix (e.g., `APP_PORT`, `APP_DATABASE_URL`)

#### Database (internal/database)

- Uses **pgxpool** for PostgreSQL connection pooling
- Uses **sqlc** for type-safe SQL query generation
- Connection pool is created at startup and shared across handlers

#### HTTP Server (internal/server)

- Uses **Gin** web framework
- Middleware:
  - Recovery (panic recovery)
  - Request logging (using slog)
- Routes are registered centrally in `server.go`

#### Handlers (internal/handler)

- Each handler receives dependencies via constructor injection
- Handlers are grouped by domain (e.g., `health.go`)

## Design Decisions

### Why Gin?

- High performance with low memory footprint
- Large ecosystem and community support
- Built-in middleware support
- Easy to learn and use

### Why pgxpool?

- Native PostgreSQL driver (not database/sql wrapper)
- Better performance than database/sql
- Built-in connection pooling
- Support for PostgreSQL-specific features

### Why sqlc?

- Type-safe SQL at compile time
- No runtime reflection overhead
- SQL stays in `.sql` files (not embedded in Go code)
- Easy to review and test queries

### Why Cobra + Viper?

- Industry standard for Go CLI applications
- Cobra provides structured CLI with subcommands
- Viper provides flexible configuration with multiple sources
- Both work well together with flag binding

### Why slog?

- Standard library (no external dependency)
- Structured logging with levels
- JSON output for production environments
- Easy to integrate with logging infrastructure

## API Design

### Health Endpoints

- `GET /health` - Liveness probe
  - Always returns 200 if server is running
  - Used by orchestrators to detect crashed containers

- `GET /ready` - Readiness probe
  - Returns 200 only if all dependencies are healthy
  - Checks database connectivity
  - Used by load balancers to route traffic

## Security Considerations

- Database credentials via environment variables (not hardcoded)
- Non-root user in Docker container
- Connection timeouts configured
- Input validation on handlers (add as needed)

## Future Considerations

- Add authentication middleware (JWT)
- Add rate limiting
- Add request validation
- Add CORS configuration
- Add metrics (Prometheus)
- Add tracing (OpenTelemetry)
