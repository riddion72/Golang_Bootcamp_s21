# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./cmd/client/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /server /app/server
COPY --from=builder /client /app/client
COPY .env .

EXPOSE 8081 8082