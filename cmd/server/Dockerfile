FROM golang:alpine3.21 as builder

WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/server ./cmd/server
COPY ./internal ./internal
RUN go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /server
RUN apk update && apk add curl
COPY --from=builder /server/server ./

CMD ["./server"]