# Resolve Service

A microservice for resolving dice rolls across various tabletop RPG systems.

## Overview

Resolve Service is the core component of DiceDasher responsible for processing dice rolls. The service supports multiple game systems through a modular resolver architecture.

## Supported Systems

| System    | Description                         | Actions                    |
|-----------|-------------------------------------|----------------------------|
| `generic` | Universal NdM dice rolls            | `roll`                     |
| `vtmv5`   | Vampire: The Masquerade 5th Edition | `roll`, `reroll`, `check`  |
| `tes`     | The Electric State                  | `roll`, `push`             |

## Architecture

```
services/resolve/
├── cmd/resolve-service/    # Entry point
├── internal/
│   ├── config/             # Configuration
│   ├── handler/            # HTTP handlers
│   ├── repository/         # Database layer
│   └── system/             # Game system resolvers
│       ├── generic/        # Universal rolls
│       ├── vtmv5/          # VTM v5
│       └── tes/            # The Electric State
```

## Running

### Locally

```bash
cd services/resolve
go run cmd/resolve-service/main.go
```

### Docker

```bash
docker build -t resolve-service -f services/resolve/Dockerfile .
docker run --env-file services/resolve/.env -p 8080:8080 resolve-service
```

## Configuration

The service is configured via environment variables:

| Variable       | Description                               | Default   |
|----------------|-------------------------------------------|-----------|
| `HTTP_ADDR`    | HTTP listen address                       | `:8080`   |
| `DATABASE_URL` | PostgreSQL connection URL                 | —         |
| `LOG_MODE`     | `default`, `debug`, or `contrast` logging | `default` |

The default logger omits `/health` requests. Use `LOG_MODE=debug` to include them.

## API

Full endpoint documentation: [endpoints.md](endpoints.md)

### Quick Start

```bash
# Health check
curl http://localhost:8080/health

# Generic roll (2d6)
curl -X POST "http://localhost:8080/resolve?system=generic" \
  -H "Content-Type: application/json" \
  -d '{"number": 2, "size": 6}'

# VTM v5 roll
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"attribute": 3, "skill": 2, "hunger": 2, "target": 3}'
```

## Adding a New System

1. Create a package in `internal/system/<system_name>/`
2. Implement the `system.Resolver` interface:
   ```go
   type Resolver interface {
       Resolve(ctx context.Context, action string, raw json.RawMessage) (any, error)
   }
   ```
3. Import the package in `main.go` and add its resolver to the map passed to
   `service.New`. Construct any dependencies there; there is no global registry.
4. Return `system.RequestError` for malformed or invalid commands. Return execution
   errors with their causes intact; the HTTP handler maps errors to responses.

## Execution and persistence

`handler.Handler` decodes HTTP requests and maps errors to HTTP statuses.
`service.ResolveService` selects a resolver, executes the command, serializes the
result and private state, and saves both in one history record. Resolvers return
`(result, error)` without HTTP status codes. TES/VTM accept a `system.HistoryReader`
through their constructors; the service accepts a `HistoryWriter`. Dependencies
are assembled in `main.go`, without storing repositories in request contexts.

Every successful operation includes a persisted `record_id`. A storage failure
returns HTTP 500 with a generic public message; the underlying error is logged.
The service does not retry the random operation. This does not provide idempotency
for client retries or guarantee delivery after a successful database write.

## Persisted roll state

TES and VTM V5 store a private `state_payload` snapshot alongside request/response history.
It contains the dice, target, and original/parent record IDs. The public response is unchanged.
The initial roll has no parent; its record ID becomes the original ID on continuation.
Allowed transitions are TES `roll|push -> push` and VTM V5 `roll|reroll -> reroll`.
This does not enforce a maximum number of continuations or prevent branching from the same record.

The `state_payload` column is defined in `database/init/001_resolve.sql`
and is created when initializing `resolve_db`.
Legacy initial rolls can still be continued. Legacy push/reroll records without a snapshot
return HTTP 409; continue from the original roll instead. Wrong system/action also returns
409, missing history returns 404, and invalid reroll indices return 422.

## Dice pool limits

All systems allow at most 1000 dice per complete roll (including all components),
and dice may have 2–1,000,000 sides. Invalid request pools return HTTP 422.
The limits are defined in `pkg/dice/pool.go`; pool arithmetic rejects integer overflow.
TES applies modifiers to the attribute pool: it must remain nonnegative, and the
combined attribute and gear pool must be nonempty. VTM hunger dice cannot exceed
`attribute + skill`. Empty component pools are allowed. Persisted continuation
states and the shared dice functions enforce the same resource limits.
