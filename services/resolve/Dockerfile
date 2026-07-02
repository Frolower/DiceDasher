FROM golang:1.26-alpine AS build

WORKDIR /app

COPY pkg/go.mod ./pkg/
COPY services/resolve/go.mod services/resolve/go.sum ./services/resolve/

WORKDIR /app/services/resolve
RUN go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/resolve/ ./services/resolve/

WORKDIR /app/services/resolve
RUN go build -o /app/resolve-service ./cmd/resolve-service

FROM alpine:3.21

RUN apk add --no-cache curl

WORKDIR /app
COPY --from=build /app/resolve-service /app/resolve-service

EXPOSE 8080
CMD ["/app/resolve-service"]
