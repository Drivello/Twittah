# --- Build Stage ---
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o tweetservice ./cmd/...

# --- Runtime Stage ---
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/tweetservice .
COPY .env-example .env
EXPOSE 8080
ENV GIN_MODE=release
# IMPORTANT: Pass real secrets/variables at runtime, never in the image.
CMD ["./tweetservice"]
