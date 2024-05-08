FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o mongodb-nats-connector ./cmd/mongodb-nats-connector/main.go
CMD ./mongodb-nats-connector
