# ELA-BE

Employee List Application backend scaffold using Go, Gin, Go Kit, GORM, PostgreSQL, Viper, and Logrus.

## Project structure

```
.
├── cmd/api                  # Application entrypoint
├── internal/config          # Viper config loading
├── internal/logger          # Logrus initialization
├── internal/database        # PostgreSQL connection (GORM)
├── internal/employee        # Employee domain model/repo/service/endpoints/handlers
├── internal/http            # Gin router setup
├── configs/config.yaml      # Default local config
├── .env.example             # Environment variable example
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Prerequisites

- Go 1.22+
- Docker + Docker Compose (optional, for local postgres/app)

## Configuration

Defaults are in `configs/config.yaml`.
Environment variables can override defaults (see `.env.example`).

## Run locally

1. Start PostgreSQL:

```bash
docker compose up -d db
```

2. Run the API:

```bash
make run
```

API health check:

```bash
curl http://localhost:8080/health
```

## Employee endpoints

Base path: `/api/v1/employees`

- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:id`
- `PUT /api/v1/employees/:id`
- `DELETE /api/v1/employees/:id`

Example create payload:

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "position": "Software Engineer"
}
```

## Run with Docker

```bash
docker compose up --build
```

## Development commands

```bash
make tidy
make fmt
make test
make build
```
