# Resolve Service

A microservice for resolving dice rolls across various tabletop RPG systems.

## Overview

Resolve Service is the core component of DiceDasher responsible for processing dice rolls. The service supports multiple game systems through a modular resolver architecture.

## Supported Systems

| System    | Description                                      | Actions                    |
|-----------|--------------------------------------------------|----------------------------|
| `generic` | Universal NdM dice rolls                         | `roll`                     |
| `vtmv5`   | Vampire: The Masquerade 5th Edition              | `roll`, `reroll`, `check`  |
| `tes`     | Tales from the Eternal Steppe (Year Zero Engine) | `roll`, `push`             |

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
│       └── tes/            # Tales from the Eternal Steppe
```

## Running

### Locally

```bash
cd services/resolver
go run cmd/resolver-service/main.go
```

### Docker

```bash
docker build -t resolver-service -f services/resolver/Dockerfile .
docker run -p 8080:8080 resolver-service
```

## Configuration

The service is configured via environment variables:

| Variable       | Description              | Default |
|----------------|--------------------------|---------|
| `PORT`         | HTTP server port         | `8080`  |
| `DATABASE_URL` | PostgreSQL connection URL | —      |

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
       Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error)
   }
   ```
3. Register the resolver in `init()`:
   ```go
   func init() {
       system.Register("system_name", Resolver{})
   }
   ```
4. Import the package in `main.go`
