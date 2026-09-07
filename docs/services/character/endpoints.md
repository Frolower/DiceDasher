# Character Service Endpoints

Base URL (local): `http://localhost:8081`

## GET /health

Health check endpoint.

**Request:** No body required.

**Response:**

- `200 OK` — service is healthy
- Body: `OK` (`text/plain`)

```bash
curl -i http://localhost:8081/health
```

## POST /character

Validate, create, and store a character. The request body depends on the selected system.

### Query parameters

| Parameter | Required | Description                           |
|-----------|----------|---------------------------------------|
| `system`  | Yes      | Game system identifier; currently `tes` |

### Responses

| Status | Description                                      |
|--------|--------------------------------------------------|
| `201`  | Character created and stored                     |
| `400`  | Missing system, malformed JSON, or missing user ID |
| `404`  | Unknown system                                   |
| `422`  | The character violates integrity or creation rules   |
| `500`  | Database or internal server error                |

### Success response

```json
{
  "id": "90e8e1f0-c213-4fa3-9a64-7d7ce3115e04"
}
```

The returned `id` is the UUID of the new row in the character database.

### Example

```bash
curl -X POST "http://localhost:8081/character?system=tes" \
  -H "Content-Type: application/json" \
  --data @character.json
```

See [systems/tes.md](systems/tes.md) for the complete request schema.

### Validation errors

A `422` response uses `application/json` and contains all detected violations:

```json
{
  "errors": [
    {"field": "character.name", "code": "required", "message": "name is required"},
    {"field": "character.gear[1].price", "code": "negative_value", "message": "must be greater or equal to 0"}
  ]
}
```

Paths use JSON field names and zero-based array indexes. Codes are `required`, `too_long`, `negative_value`, `invalid_value`, `invalid_item`, `incompatible_field`, and `creation_rule`. Some type-specific and creation checks identify a group such as `character.gear[0]` or `character.stats`; their message describes the constraint. Other error statuses retain their existing text format.

Compatibility: free-mode requests with missing/blank names or invalid supplied values now return `422`; validation errors previously returned plain text.
