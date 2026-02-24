FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
# Note: CGO is disabled for pure Go sqlite driver compatibility (modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o xxxclaw ./cmd/xxxclaw

# Final minimal image
FROM scratch

# Copy timezone data & certs
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

# Create data directory for SQLite
VOLUME /app/data

# Copy the pre-built binary
COPY --from=builder /app/xxxclaw .

# Default ports: HTTP config map
EXPOSE 8080 6060

ENTRYPOINT ["/app/xxxclaw"]
