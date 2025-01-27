# Build Stage
FROM golang:1.23 as builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy scripts
COPY start.sh .
COPY wait-for-it.sh .
RUN chmod +x start.sh wait-for-it.sh

# Copy the source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o callback cmd/callback/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o worker cmd/worker/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o worker_manager cmd/worker_manager/main.go

# Runtime Stage
# Runtime Stage
FROM alpine:latest

WORKDIR /root/

# Install dependencies
RUN apk add --no-cache bash ca-certificates

# Copy binaries and scripts from the builder stage
COPY --from=builder /app/start.sh .
COPY --from=builder /app/wait-for-it.sh .
COPY --from=builder /app/callback .
COPY --from=builder /app/worker .
COPY --from=builder /app/worker_manager .

# Copy static assets (e.g., HTML files)
COPY --from=builder /app/static ./static

# Make scripts executable
RUN chmod +x start.sh wait-for-it.sh

# Expose required port
EXPOSE 9090

# Use the start.sh script to start the service
CMD ["./start.sh"]
