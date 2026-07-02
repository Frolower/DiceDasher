# API Reference

This document provides a general list of all API endpoints across DiceDasher services.

Interactive API documentation is available through Swagger UI at `http://localhost:8082` after running `make up`. The source specification is [openapi.yaml](openapi.yaml).

## Resolve Service

Base URL: `http://localhost:8080`

Full documentation: [/docs/services/resolve/](../services/resolve/)

### Endpoints

| Method | Endpoint   | Description       | Details                                            |
|--------|------------|-------------------|----------------------------------------------------|
| GET    | `/health`  | Health check      | [→](../services/resolve/endpoints.md#get-health)   |
| POST   | `/resolve` | Resolve dice roll | [→](../services/resolve/endpoints.md#post-resolve) |

### Supported Systems

| System    | Actions                   | Description                         |
|-----------|---------------------------|-------------------------------------|
| `generic` | `roll`                    | Universal NdM dice rolls            |
| `vtmv5`   | `roll`, `reroll`, `check` | Vampire: The Masquerade 5th Edition |
| `tes`     | `roll`, `push`            | The Electric State                  |

### Quick Examples

```sh
# Health check
curl http://localhost:8080/health

# Generic: Roll 2d8
curl -X POST "http://localhost:8080/resolve?system=generic" \
  -H "Content-Type: application/json" \
  -d '{"number": 2, "size": 8}'

# VTM v5: Dice pool roll
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"attribute": 3, "skill": 2, "hunger": 2, "target": 3}'

# TES: Attribute + gear roll
curl -X POST "http://localhost:8080/resolve?system=tes&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"attr": 4, "assist": 1, "gear": 2, "modificator": 0, "target": 2}'
```

## Character Service

Base URL: `http://localhost:8081`

Full documentation: [/docs/services/character/](../services/character/)

### Endpoints

| Method | Endpoint     | Description                  | Details                                               |
|--------|--------------|------------------------------|-------------------------------------------------------|
| GET    | `/health`    | Health check                 | [→](../services/character/endpoints.md#get-health)    |
| POST   | `/character` | Create and store a character | [→](../services/character/endpoints.md#post-character) |

### Supported Systems

| System | Operation        | Description        |
|--------|------------------|--------------------|
| `tes`  | Character create | The Electric State |

### Quick Example

```sh
curl -X POST "http://localhost:8081/character?system=tes" \
  -H "Content-Type: application/json" \
  --data @character.json
```
