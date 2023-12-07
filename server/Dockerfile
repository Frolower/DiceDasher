FROM golang:1.20.5-alpine AS builder

WORKDIR /usr/local/src

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./main.go


FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /usr/local/bin/app

# connections to front and BD
EXPOSE 8080
EXPOSE 3232

CMD ["app"]