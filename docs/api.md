# ELA-BE API Documentation

Employee List Application — REST API reference.

- **Base URL (local):** `http://localhost:8080`
- **API version prefix:** `/api/v1`
- **Content-Type:** `application/json`

---

## Table of contents

1. [Health check](#health-check)
2. [Employee endpoints](#employee-endpoints)
   - [Create an employee](#create-an-employee)
   - [List all employees](#list-all-employees)
   - [Get an employee by ID](#get-an-employee-by-id)
   - [Update an employee](#update-an-employee)
   - [Delete an employee](#delete-an-employee)
3. [Employee object](#employee-object)
4. [Error responses](#error-responses)
5. [Status code summary](#status-code-summary)

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
| `200` | OK               | GET endpoints, PUT                   |
| `201` | Created          | POST (create employee)               |
| `204` | No Content       | DELETE                               |
| `400` | Bad Request      | Validation / parse errors            |
| `404` | Not Found        | Resource does not exist              |
| `500` | Internal Server Error | Unexpected / unhandled errors   |
