# -----------------------------
# Stage 1: Build the Go binary
# -----------------------------
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN go build -o api-service main.go

# -----------------------------
# Stage 2: Minimal runtime image
# -----------------------------
FROM gcr.io/distroless/base-debian12

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/api-service .

# Expose port used by your API
EXPOSE 8080

# Run the binary
CMD ["./api-service"]

