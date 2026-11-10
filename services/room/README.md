# Room service

Minimal gRPC server. Room RPCs and persistence will be added when the API is defined.

## Run

From the repository root:

```sh
go run ./services/room/cmd/room-service
```

`GRPC_ADDR` sets the listen address (default `:50051`). The server currently uses
plaintext gRPC for local development and trusted internal connections.

## Health check

The standard `grpc.health.v1.Health` service provides `Check` and `Watch`.
Use an empty service name to check the whole server. It reports `SERVING` while
running and transitions to `NOT_SERVING` when shutdown begins.

Reflection is enabled, so `grpcurl` works without local proto files:

```sh
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"service":""}' localhost:50051 grpc.health.v1.Health/Check
```

Expected health response:

```json
{"status":"SERVING"}
```

Through the backend: `GET http://localhost:8000/health?service=room` returns HTTP
200 when room reports `SERVING`, or 503 when it is unhealthy or unreachable.
Room is also included in Docker Compose (`make up`).

SIGINT/SIGTERM initiates graceful shutdown, with a 10-second limit before forcing
active RPCs to stop. No separate HTTP health endpoint is needed.

## Extend

Register generated room RPC implementations in `internal/server.New`.
The built-in health service already includes its proto definitions and Go code.
