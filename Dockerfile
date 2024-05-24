FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o mongodb-nats-connector ./cmd/connector/main.go
CMD ./mongodb-nats-connector
