FROM golang:1.26-alpine AS build

WORKDIR /app

COPY pkg/go.mod ./pkg/
COPY services/resolver/go.mod services/resolver/go.sum ./services/resolver/

WORKDIR /app/services/resolver
RUN go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/resolver/ ./services/resolver/

WORKDIR /app/services/resolver
RUN go build -o /app/resolver-service ./cmd/resolver-service

FROM alpine:3.21

WORKDIR /app
COPY --from=build /app/resolver-service /app/resolver-service

EXPOSE 8080
CMD ["/app/resolver-service"]
