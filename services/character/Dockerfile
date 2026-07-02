FROM golang:1.26-alpine AS build

WORKDIR /app

COPY pkg/go.mod ./pkg/
COPY services/character/go.mod services/character/go.sum ./services/character/

WORKDIR /app/services/character
RUN go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/character/ ./services/character/

WORKDIR /app/services/character
RUN go build -o /app/character-service ./cmd/character-service

FROM alpine:3.21

RUN apk add --no-cache curl

WORKDIR /app
COPY --from=build /app/character-service /app/character-service

EXPOSE 8081
CMD ["/app/character-service"]
