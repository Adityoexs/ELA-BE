# ELA-BE API Documentation

Employee List Application — REST API reference.

- **Base URL (local):** `http://localhost:8080`
- **API version prefix:** `/api/v1`
- **Content-Type:** `application/json`

---

## Table of contents

1. [Authentication](#authentication)
   - [Login](#login)
   - [Authenticated usage](#authenticated-usage)
2. [Health check](#health-check)
3. [Employee endpoints](#employee-endpoints)
   - [Create an employee](#create-an-employee)
   - [List all employees](#list-all-employees)
   - [Get an employee by ID](#get-an-employee-by-id)
   - [Update an employee](#update-an-employee)
   - [Delete an employee](#delete-an-employee)
4. [Employee object](#employee-object)
5. [Error responses](#error-responses)
6. [Status code summary](#status-code-summary)

---

## Authentication

All `/api/v1/employees` endpoints require a valid **Bearer JWT** in the `Authorization` header.

### Login

```
POST /api/v1/auth/login
```

Returns a signed JWT on successful authentication.

**Request body**

| Field      | Type   | Required | Description                |
|------------|--------|----------|----------------------------|
| `email`    | string | ✅ yes   | Registered email address   |
| `password` | string | ✅ yes   | Password for the account   |

```json
{
  "email": "admin@example.com",
  "password": "admin123"
}
```

**Response `200 OK`**

```json
{
  "token": "<signed JWT>"
}
```

**Error responses**

| Status | Condition                              | Body                                      |
|--------|----------------------------------------|-------------------------------------------|
| `400`  | Missing or invalid `email`/`password`  | `{"error": "..."}`                        |
| `401`  | Wrong email or password                | `{"error": "invalid email or password"}`  |

**Example curl**

```bash
# Obtain a JWT
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

> **MVP note:** The default user (`admin@example.com` / `admin123`) is hardcoded in
> `internal/auth/service.go`. Replace the `mvpUsers` map with a real database lookup
> once a users table is introduced.

---

### Authenticated usage

Pass the token returned by `POST /api/v1/auth/login` as a `Bearer` token in every
request to a protected endpoint:

```bash
# Store the token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' | jq -r .token)

# Use it to list employees
curl -s http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer $TOKEN"
```

**401 responses returned when the token is:**

| Scenario              | Response body                                            |
|-----------------------|----------------------------------------------------------|
| Header absent         | `{"error": "authorization header required"}`             |
| Non-Bearer scheme     | `{"error": "authorization header must use Bearer scheme"}` |
| Invalid / expired     | `{"error": "invalid or expired token"}`                  |

**JWT configuration** (via `configs/config.yaml` or environment variables):

| Config key              | Env variable          | Default                        |
|-------------------------|-----------------------|--------------------------------|
| `jwt.secret`            | `JWT_SECRET`          | `change-me-in-production`      |
| `jwt.expiry_seconds`    | `JWT_EXPIRY_SECONDS`  | `3600` (1 hour)                |

> **Security note:** Always set `JWT_SECRET` to a strong random value in production.

---

## Health check

```
GET /health
```

Returns the health status of the running service. No authentication required.

**Response `200 OK`**

```json
{
  "status": "ok"
}
```

---

## Employee endpoints

All employee endpoints share the base path `/api/v1/employees`.

### Create an employee

```
POST /api/v1/employees
```

Creates a new employee record.

**Request body**

| Field      | Type   | Required | Description                          |
|------------|--------|----------|--------------------------------------|
| `name`     | string | ✅ yes   | Full name of the employee            |
| `email`    | string | ✅ yes   | Unique email address of the employee |
| `position` | string | ❌ no    | Job title / position                 |

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "position": "Software Engineer"
}
```

**Response `201 Created`**

Returns the newly created [Employee object](#employee-object).

```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "position": "Software Engineer",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error responses**

| Status | Condition                                    | Body                                      |
|--------|----------------------------------------------|-------------------------------------------|
| `400`  | Missing `name` or `email`, malformed JSON    | `{"error": "name is required"}`           |
| `400`  | Duplicate `email` (database constraint)      | `{"error": "create employee: ..."}`       |
| `500`  | Unexpected server error                      | `{"error": "internal server error"}`      |

---

### List all employees

```
GET /api/v1/employees
```

Returns a list of all employees, ordered by `id` ascending.

**Request body:** none

**Response `200 OK`**

Returns a JSON array of [Employee objects](#employee-object). Returns an empty array `[]` when there are no employees.

```json
[
  {
    "id": 1,
    "name": "Jane Doe",
    "email": "jane@example.com",
    "position": "Software Engineer",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "name": "John Smith",
    "email": "john@example.com",
    "position": "Product Manager",
    "created_at": "2024-01-16T09:00:00Z",
    "updated_at": "2024-01-16T09:00:00Z"
  }
]
```

**Error responses**

| Status | Condition               | Body                                 |
|--------|-------------------------|--------------------------------------|
| `500`  | Unexpected server error | `{"error": "internal server error"}` |

---

### Get an employee by ID

```
GET /api/v1/employees/:id
```

Returns a single employee by their numeric ID.

**Path parameters**

| Parameter | Type    | Description                 |
|-----------|---------|-----------------------------|
| `id`      | integer | Unique identifier (uint > 0)|

**Request body:** none

**Response `200 OK`**

Returns the [Employee object](#employee-object).

```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "position": "Software Engineer",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error responses**

| Status | Condition                        | Body                                 |
|--------|----------------------------------|--------------------------------------|
| `400`  | `id` is not a valid integer      | `{"error": "invalid employee id"}`   |
| `404`  | Employee not found               | `{"error": "employee not found"}`    |
| `500`  | Unexpected server error          | `{"error": "internal server error"}` |

---

### Update an employee

```
PUT /api/v1/employees/:id
```

Updates an existing employee. Only the fields provided in the request body are applied — omitted or blank fields are left unchanged (partial update).

**Path parameters**

| Parameter | Type    | Description                  |
|-----------|---------|------------------------------|
| `id`      | integer | Unique identifier (uint > 0) |

**Request body**

All fields are optional. At least one non-blank value should be provided for the update to have any effect.

| Field      | Type   | Required | Description                    |
|------------|--------|----------|--------------------------------|
| `name`     | string | ❌ no   | New full name                  |
| `email`    | string | ❌ no   | New email address (must be unique) |
| `position` | string | ❌ no   | New job title / position       |

```json
{
  "position": "Senior Software Engineer"
}
```

**Response `200 OK`**

Returns the updated [Employee object](#employee-object).

```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "position": "Senior Software Engineer",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-20T14:00:00Z"
}
```

**Error responses**

| Status | Condition                        | Body                                 |
|--------|----------------------------------|--------------------------------------|
| `400`  | `id` is not a valid integer      | `{"error": "invalid employee id"}`   |
| `400`  | Malformed JSON body              | `{"error": "..."}`                   |
| `400`  | Duplicate `email` constraint     | `{"error": "update employee: ..."}`  |
| `404`  | Employee not found               | `{"error": "employee not found"}`    |
| `500`  | Unexpected server error          | `{"error": "internal server error"}` |

---

### Delete an employee

```
DELETE /api/v1/employees/:id
```

Permanently deletes an employee record.

**Path parameters**

| Parameter | Type    | Description                  |
|-----------|---------|------------------------------|
| `id`      | integer | Unique identifier (uint > 0) |

**Request body:** none

**Response `204 No Content`**

Empty body on success.

**Error responses**

| Status | Condition                        | Body                                 |
|--------|----------------------------------|--------------------------------------|
| `400`  | `id` is not a valid integer      | `{"error": "invalid employee id"}`   |
| `404`  | Employee not found               | `{"error": "employee not found"}`    |
| `500`  | Unexpected server error          | `{"error": "internal server error"}` |

---

## Employee object

The canonical representation of an employee returned by the API.

| Field        | Type             | Description                                  |
|--------------|------------------|----------------------------------------------|
| `id`         | integer          | Auto-generated primary key                   |
| `name`       | string (≤ 255)   | Full name — required, non-blank              |
| `email`      | string (≤ 255)   | Email address — required, unique             |
| `position`   | string (≤ 255)   | Job title — optional, may be empty string    |
| `created_at` | string (RFC 3339)| Timestamp when the record was created (UTC)  |
| `updated_at` | string (RFC 3339)| Timestamp when the record was last updated (UTC) |

---

## Error responses

All error responses share the same JSON shape:

```json
{
  "error": "<human-readable description>"
}
```

---

## Status code summary

| Code  | Meaning          | Used by                              |
|-------|------------------|--------------------------------------|
| `200` | OK               | GET endpoints, PUT, Login            |
| `201` | Created          | POST (create employee)               |
| `204` | No Content       | DELETE                               |
| `400` | Bad Request      | Validation / parse errors            |
| `401` | Unauthorized     | Missing, invalid, or expired JWT     |
| `404` | Not Found        | Resource does not exist              |
| `500` | Internal Server Error | Unexpected / unhandled errors   |
