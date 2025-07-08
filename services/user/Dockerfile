# ---- Build stage ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user_service ./cmd/main.go

# ---- Run stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user_service ./user_service
EXPOSE 8082
ENV USER_SERVICE_PORT=8082
ENTRYPOINT ["/app/user_service"]
