# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for a static binary
# -o /app/server builds the output to a specific file
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server .

# Copy the environment file
# Note: In a real production scenario, you would inject secrets
# securely, not copy the .env file.
COPY app.env .

# Expose the port the app runs on
EXPOSE 3334

# Command to run the executable
CMD ["/app/server"]
