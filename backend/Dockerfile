# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Set env vars for optimization
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with ldflags to reduce size
RUN go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
