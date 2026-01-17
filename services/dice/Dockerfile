FROM golang:1.25-alpine

WORKDIR /app

COPY go.work ./
COPY pkg/go.mod ./pkg/
COPY services/dice/go.mod ./services/dice/

RUN go work sync
RUN go mod download

COPY . .

RUN go build -o /app/dice-service ./services/dice/cmd/dice-service

EXPOSE 8080
CMD ["/app/dice-service"]
