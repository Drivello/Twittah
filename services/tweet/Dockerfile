FROM golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod download
RUN go build -o tweet-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tweet-service .
COPY .env .
EXPOSE 8080
CMD ["./tweet-service"]
