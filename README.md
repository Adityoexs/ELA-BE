# ELA-BE

Employee List Application backend scaffold using Go, Gin, Go Kit, GORM, PostgreSQL, Viper, and Logrus.

## Project structure

```
.
├── cmd/api                  # Application entrypoint
├── internal/auth            # JWT authentication service, middleware, handler
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

## Authentication

All `/api/v1/employees` endpoints require a valid Bearer JWT.

**Login to obtain a token:**

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

**Use the token:**

```bash
TOKEN=<paste token here>
curl http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer $TOKEN"
```

> See [`docs/api.md`](docs/api.md) for full authentication documentation.

## Employee endpoints

Base path: `/api/v1/employees`

| Method   | Path                       | Description           |
|----------|----------------------------|-----------------------|
| `POST`   | `/api/v1/employees`        | Create an employee    |
| `GET`    | `/api/v1/employees`        | List all employees    |
| `GET`    | `/api/v1/employees/:id`    | Get employee by ID    |
| `PUT`    | `/api/v1/employees/:id`    | Update an employee    |
| `DELETE` | `/api/v1/employees/:id`    | Delete an employee    |

Example create payload:

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "position": "Software Engineer"
}
```

📄 **Full API reference:** [`docs/api.md`](docs/api.md)  
📐 **OpenAPI 3.0 spec:** [`docs/openapi.yaml`](docs/openapi.yaml)

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
