FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /app/dice-service ./cmd/dice-service

EXPOSE 8080

RUN useradd app
USER app

CMD ["/app/dice-service"]
