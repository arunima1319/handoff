# Stage 1: Build the Go application
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/server .

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/app ./app

EXPOSE 8080

CMD ["./server"]