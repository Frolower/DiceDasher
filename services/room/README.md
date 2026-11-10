# Room service

Minimal gRPC server with health checking and a CreateRoom RPC stub.

## CreateRoom round trip

Send an authenticated request to the backend:

```sh
curl -i http://localhost:8000/room \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"system":"vtmv5"}'
```

The backend takes the user UUID from the verified access token and calls
`dicedasher.room.v1.RoomService/CreateRoom` with `system` and `user_id`.
The room handler logs both fields and returns `status: "ok"` and a random UUID v4
in `room_id`; the backend forwards that UUID in its HTTP 200 response:
`{"status":"ok","room_id":"9da3cc51-c9a2-4f98-a865-d4aa135c0e96"}`.
Every call generates a new UUID. No room is created or persisted yet.
The `system` must be nonblank; there is no list of allowed systems yet.
The HTTP body accepts only `system`; a client-supplied `user_id` is rejected.

Missing/invalid authorization returns 401, invalid input returns 400, oversized
bodies return 413, unavailable room returns 503, a timeout returns 504, and other
upstream errors return 502. Backend RPC calls have a 3-second timeout and inherit
HTTP request cancellation. `ROOM_GRPC_ADDR` configures the backend's room address
(default `localhost:50051`; Docker Compose uses `room:50051`).

The shared contract and generated clients are in `pkg/room/v1`.

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
