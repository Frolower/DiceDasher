FROM golang:1.27-alpine AS build
WORKDIR /app

COPY pkg/ ./pkg/
COPY services/backend/go.mod services/backend/go.sum ./services/backend/
COPY services/backend/cmd/ ./services/backend/cmd/
COPY services/backend/internal/ ./services/backend/internal/
RUN go work init ./pkg ./services/backend
RUN CGO_ENABLED=0 go build -o /app/backend-service ./services/backend/cmd/backend-service

FROM alpine:3.21
RUN apk add --no-cache curl
WORKDIR /app
COPY --from=build /app/backend-service /app/backend-service
EXPOSE 8083
CMD ["/app/backend-service"]
