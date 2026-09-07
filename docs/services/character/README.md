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
    ├── handler/           # HTTP parsing and error/status mapping
    ├── service/           # Creation use case and repository interface
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
       CreateCharacter(ctx context.Context, raw json.RawMessage) (CreatedCharacter, error)
   }
   ```

3. Pass the implementation to the service constructor in `cmd/character-service/main.go`:

   ```go
   characters := service.New(repository.New(repo.Pool()), map[string]system.Character{
       "tes": tes.Character{},
       "system_name": implementation,
   })
   handler.New(characters).RegisterRouters(r)
   ```

## Creation flow and dependencies

`Handler → CharacterService → system.Character → Repository`

The service selects a strategy, prepares the character, and persists it through its `Repository` interface. `main` constructs the PostgreSQL adapter, service, and handler explicitly. The service snapshots its system map; there is no global registration or repository stored in request context. Context carries request cancellation and logging context through the use case.

Strategies return a character or an error without HTTP status codes. The handler maps `system.InputError` to 400, `system.ValidationError` to the existing 422 JSON response, `service.ErrUnknownSystem` to 404, and persistence/internal failures to 500 without exposing their causes. Success remains 201 with an `id`. Request syntax is checked before invoking the service, so malformed JSON with an unknown system returns 400.

Service tests use an in-memory repository stub; HTTP tests exercise routing and response mapping without PostgreSQL.
