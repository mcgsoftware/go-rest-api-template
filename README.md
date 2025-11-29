# go-rest-api-template

A production-ready Go REST API template with Gin, PostgreSQL, and modern tooling.

## Features

- **Gin** web framework
- **PostgreSQL** with pgxpool connection pooling
- **sqlc** for type-safe SQL queries
- **Cobra + Viper** for CLI and configuration
- **slog** structured logging
- **OpenAPI/Swagger** documentation
- **Docker + Compose** for containerization
- **Makefile** for common tasks

## Quick Start

### Prerequisites

- Go 1.23+
- PostgreSQL
- Docker (optional)

### Install Dependencies

```bash
go mod download
```

### Install Tools

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### Run with Docker

```bash
# Start PostgreSQL and API
make docker-up

# Stop containers
make docker-down
```

### Run Locally

```bash
# Generate code
make generate

# Build
make build

# Run
./bin/server serve --port 8080 --db "postgres://user:pass@localhost:5432/mydb"
```

## Configuration

Configuration can be set via CLI flags or environment variables.

| Flag | Env Variable | Default | Description |
|------|--------------|---------|-------------|
| `--port` | `APP_PORT` | 8080 | Server port |
| `--debug` | `APP_DEBUG` | false | Enable debug logging |
| `--db` | `APP_DATABASE_URL` | (required) | PostgreSQL connection URL |

**Precedence:** CLI flag > environment variable > default

### Examples

```bash
# Using flags
./bin/server serve --port 3000 --debug --db "postgres://localhost/mydb"

# Using environment variables
APP_PORT=3000 APP_DEBUG=true APP_DATABASE_URL="postgres://localhost/mydb" ./bin/server serve

# Using .env file (copy from .env.example)
cp .env.example .env
# Edit .env with your values
./bin/server serve
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness check (always 200) |
| GET | `/ready` | Readiness check with DB connectivity |
| GET | `/swagger/*` | Swagger UI |

## Makefile Commands

```bash
make build         # Build the binary
make run           # Run locally
make test          # Run tests
make docker-build  # Build Docker image
make docker-up     # Start with docker-compose
make docker-down   # Stop containers
make swagger       # Generate Swagger docs
make sqlc          # Generate sqlc code
make generate      # Run all code generation (sqlc + swagger)
make help          # Show all available commands
```

## sqlc Usage

### Writing Queries

Add SQL queries to `db/queries/*.sql`:

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;

-- name: SearchUsersByDateRange :many
-- Call a PostgreSQL function that returns multiple rows
SELECT * FROM search_users_by_date_range($1::timestamp, $2::timestamp);
```

### Shell Commands

```bash
# Generate Go code from SQL queries
sqlc generate -f db/sqlc.yaml

# Verify queries without generating
sqlc compile -f db/sqlc.yaml

# Check for errors
sqlc vet -f db/sqlc.yaml
```

### Makefile Commands

```bash
# Generate sqlc code
make sqlc

# Generate all code (sqlc + swagger)
make generate

# Build after generating
make generate build

# Full dev workflow
make generate build run
```

### Example Go Usage

```go
// Call PostgreSQL function with 2 inputs, returns multiple rows
users, err := queries.SearchUsersByDateRange(ctx, db.SearchUsersByDateRangeParams{
    Column1: startDate,  // $1::timestamp
    Column2: endDate,    // $2::timestamp
})
if err != nil {
    return err
}
for _, user := range users {
    fmt.Printf("User: %s\n", user.Name)
}
```

## Project Structure

```
.
├── cmd/server/main.go           # Entry point (Cobra CLI)
├── internal/
│   ├── config/config.go         # Viper configuration
│   ├── database/postgres.go     # pgxpool connection
│   ├── handler/health.go        # Health check handlers
│   └── server/server.go         # Gin server setup
├── api/openapi.yaml             # OpenAPI specification
├── db/
│   ├── schema.sql               # PostgreSQL schema
│   ├── sample-data.sql          # Example seed data
│   ├── sqlc.yaml                # sqlc configuration
│   └── queries/                 # SQL query files
├── doc/design.md                # Design documentation
├── Dockerfile
├── docker-compose.yaml
├── Makefile
└── .env.example
```

## Go Dependencies

```bash
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/files
```

## License

MIT
