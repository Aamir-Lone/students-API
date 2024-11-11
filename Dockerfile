# Build stage
FROM golang:1.22.5 AS builder  

# Set up the Go environment
WORKDIR /app

# Copy go.mod and go.sum files first
COPY go.mod go.sum ./

# Run go mod tidy to download dependencies
RUN go mod tidy

# Copy the entire project to the container
COPY . .


# Build the Go application
RUN go build -o students-api ./cmd/students-api

# Final stage
FROM debian:bookworm-slim  
# Use Debian 12 for runtime (includes GLIBC 2.34)

# Install the necessary dependencies
RUN apt-get update && apt-get install -y \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

# Set up the application directory
WORKDIR /app

# Copy the compiled Go binary from the builder image
COPY --from=builder /app/students-api /app/

# Set the entry point to run the Go binary
CMD ["./students-api"]

# Expose the application port
EXPOSE 8080
