# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12

# Copy the built static binary from the builder stage
COPY --from=builder /server /server

# Expose port (metadata)
EXPOSE 8080

# Command to run
CMD ["/server"]
