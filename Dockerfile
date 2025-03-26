# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.21
# Create directory for secrets
RUN mkdir -p /secrets
# Copy the binary
COPY --from=builder /server /server
# Create non-root user
RUN adduser -D appuser
USER appuser

# Expose port
EXPOSE 8080

# Command to run
CMD ["/server"]
