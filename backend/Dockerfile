# ============================================
# STAGE 1: Builder
# ============================================
FROM golang:1.26-alpine AS builder
# Install dependencies for build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
# Copy go.mod and go.sum first (leverages Docker cache)
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy source code
COPY . .
# Compile binary (CGO_ENABLED=0 for smaller final image)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api
# ============================================
# STAGE 2: Production
# ============================================
FROM alpine:3.19
WORKDIR /root/
# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates
# Copy compiled binary from builder
COPY --from=builder /app/server .
# Environment variables for configuration
ENV PORT=8080
# Expose port
EXPOSE 8080
# Non-root user for security
RUN adduser -D -g '' appuser
USER appuser
CMD ["./server"]