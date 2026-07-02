# Character Service

A microservice for validating, creating, and storing characters for supported tabletop RPG systems.

## Overview

Character Service dispatches each request to a system-specific character implementation. It currently supports character creation for The Electric State (`tes`). Created characters are stored in PostgreSQL and identified by UUID.

## Supported Systems

| System | Description        | Operations |
|--------|--------------------|------------|
| `tes`  | The Electric State | `create`   |

## Architecture

```text
services/character/
├── cmd/character-service/ # Entry point
└── internal/
    ├── config/            # Configuration
    ├── handler/           # HTTP handlers
    ├── repository/        # Database layer
    └── system/            # System-specific character logic
        └── tes/           # The Electric State
```

## Running

### Locally

```bash
cd services/character
go run cmd/character-service/main.go
```

### Docker

```bash
docker build -t character-service -f services/character/Dockerfile .
docker run --env-file services/character/.env -p 8081:8081 character-service
```

## Configuration

| Variable       | Description                               | Default   |
|----------------|-------------------------------------------|-----------|
| `HTTP_ADDR`    | HTTP listen address                       | `:8081`   |
| `DATABASE_URL` | PostgreSQL connection URL                 | —         |
| `LOG_MODE`     | `default`, `debug`, or `contrast` logging | `default` |

`DATABASE_URL` is required.

The default logger omits `/health` requests. Use `LOG_MODE=debug` to include them.

## API

See [endpoints.md](endpoints.md) for the HTTP contract and [systems/tes.md](systems/tes.md) for the TES request model and validation rules.

### Quick Start

```bash
curl http://localhost:8081/health

curl -X POST "http://localhost:8081/character?system=tes" \
  -H "Content-Type: application/json" \
  --data @character.json
```

## Adding a New System

1. Create a package in `internal/system/<system_name>/`.
2. Implement the `system.Character` interface:

   ```go
   type Character interface {
       CreateCharacter(ctx context.Context, raw json.RawMessage) (CreatedCharacter, int, error)
   }
   ```

3. Register the implementation in `cmd/character-service/main.go`:

   ```go
   system.Register("system_name", implementation)
   ```
