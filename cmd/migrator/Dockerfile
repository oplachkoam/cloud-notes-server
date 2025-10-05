FROM golang:alpine3.21 as builder

WORKDIR /migrator
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/migrator ./cmd/migrator
COPY ./internal/migrator ./internal/migrator
RUN go build -o migrator ./cmd/migrator

FROM alpine:latest

WORKDIR /migrator
COPY --from=builder /migrator/migrator ./
COPY ./migrations ./migrations

CMD ["./migrator"]