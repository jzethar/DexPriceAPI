FROM golang:1.21 AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go mod tidy

WORKDIR /app/server
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server/main .
