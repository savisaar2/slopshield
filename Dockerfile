# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all three tools
RUN CGO_ENABLED=0 GOOS=linux go build -o slopshield ./cmd/slopshield/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o slop-prober ./cmd/slop-prober/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o slop-hunter ./cmd/slop-hunter/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the binaries from the builder stage
COPY --from=builder /app/slopshield .
COPY --from=builder /app/slop-prober .
COPY --from=builder /app/slop-hunter .

# Copy the base registry files
COPY --from=builder /app/registry ./registry

# Default to scan the current directory
ENTRYPOINT ["/app/slopshield"]
CMD ["scan", "."]
