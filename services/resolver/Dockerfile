FROM golang:1.25-alpine

WORKDIR /app

COPY go.work ./
COPY pkg/go.mod ./pkg/
COPY services/resolve/go.mod ./services/resolve/

RUN go work sync
RUN go mod download

COPY . .

RUN go build -o /app/resolver-service ./services/resolver/cmd/resolver-service

EXPOSE 8080
CMD ["/app/resolve-service"]
