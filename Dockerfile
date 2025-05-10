# Build stage
FROM golang:1.19-alpine AS builder

# Install git and SSL ca certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# Create a non-root user and group
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/app .

# Final stage
FROM alpine:3.17

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/app /app

# Set working directory
WORKDIR /

# Use the non-root user
USER appuser:appuser

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app"] 
