# Resolve Service Endpoints

Base URL (local): `http://localhost:8080`

---

## Health

### GET /health

Health check endpoint.

**Request:** No body required.

**Response:**
- `200 OK` — service is healthy
  - Body: `OK` (text/plain)

**Example:**

```bash
curl -i http://localhost:8080/health
```

---

## Resolve

### POST /resolve

Dispatcher endpoint for dice resolution. The request/response body depends on `system` + `action`.

**Query Parameters:**

| Parameter | Required | Default | Description                            |
|-----------|----------|---------|----------------------------------------|
| `system`  | Yes      | —       | Game system identifier                 |
| `action`  | No       | `roll`  | Action to perform (system-specific)    |

**Common Responses:**

| Status | Description                              |
|--------|------------------------------------------|
| 200    | Success                                  |
| 400    | Bad request (missing params, invalid JSON) |
| 404    | Unknown system                           |
| 422    | Validation error (invalid field values)  |

---

## Supported Systems

| System    | Actions                   | Description                         | Docs |
|-----------|---------------------------|-------------------------------------|------|
| `generic` | `roll`                    | Universal NdM dice rolls            | [→](systems/generic.md) |
| `vtmv5`   | `roll`, `reroll`, `check` | Vampire: The Masquerade 5th Edition | [→](systems/vtmv5.md) |
| `tes`     | `roll`, `push`            | The Electric State                  | [→](systems/tes.md) |
