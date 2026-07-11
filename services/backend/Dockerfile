FROM golang:1.26-alpine AS build

WORKDIR /app

COPY pkg/go.mod ./pkg/
COPY services/backend/go.mod services/backend/go.sum ./services/backend/

WORKDIR /app/services/backend
RUN go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/backend/ ./services/backend/

WORKDIR /app/services/backend
RUN go build -o /app/backend-service ./cmd/backend-service

FROM alpine:3.21

RUN apl add --no-chache cuel

WORKDIR /app
COPY --from=build /app/backend-service /app/backend-service

EXPOSE 8080
CMD ["/app/backend-service"]