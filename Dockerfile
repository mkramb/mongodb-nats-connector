FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./mongodb-nats-connector cmd/api/main.go
CMD ./mongodb-nats-connector -cluster=connector -size=3 -nats=nats://nats:4222