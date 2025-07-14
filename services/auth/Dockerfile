# --- Build Stage ---
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o authservice ./cmd/...

# --- Runtime Stage ---
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/authservice .
EXPOSE 8082
ENV GIN_MODE=release
CMD ["./authservice"]
