FROM golang:1.27.2-alpine AS build
WORKDIR /app
COPY pkg/ ./pkg/
COPY services/room/ ./services/room/
RUN go work init ./pkg ./services/room
RUN CGO_ENABLED=0 go build -o /app/room-service ./services/room/cmd/room-service

FROM scratch
COPY --from=build /app/room-service /room-service
EXPOSE 50051
ENTRYPOINT ["/room-service"]
