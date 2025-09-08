# Start from the official Golang base image
FROM golang:1.23 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o go-iam main.go

# -- Release Stage --
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy binary from the builder stage
COPY --from=builder /app/go-iam .
COPY --chown=nonroot:nonroot --from=builder /app/docs /docs

# Expose application port (change this if needed)
EXPOSE 3000

# Set environment variables (optional defaults)
# These can be overridden at runtime
ENV SERVER_PORT=3000

# Run the Go IAM app
ENTRYPOINT ["/go-iam"]
