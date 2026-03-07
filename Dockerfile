# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o mydns ./cmd/mydns/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/mydns .

# Expose DNS ports
EXPOSE 53/udp
EXPOSE 53/tcp

# Run the binary
CMD ["./mydns"]
