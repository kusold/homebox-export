# Stage 1: Builder
FROM golang:1.23-alpine AS builder
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build with version information and optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s \
    -X main.version=${VERSION} \
    -X main.commit=${COMMIT} \
    -X main.buildDate=${DATE}" \
    -o /app/bin/homebox-export ./cmd/homebox-export

# Stage 2: Runtime
FROM alpine:3.21 AS runtime
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
# Add basic security through non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install necessary runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Create necessary directories with proper permissions
RUN mkdir /app /export \
    && chown -R appuser:appgroup /app /export

# Copy binary from builder
COPY --from=builder --chown=appuser:appgroup /app/bin/homebox-export /app/homebox-export

# Switch to non-root user
USER appuser:appgroup

# Set working directory
WORKDIR /app

# Default export directory
VOLUME ["/export"]

# Environment variables with defaults
ENV HOMEBOX_SERVER="" \
    HOMEBOX_USER="" \
    HOMEBOX_PASS="" \
    HOMEBOX_OUTPUT="/export" \
    HOMEBOX_PAGESIZE="100"

# Set entrypoint and default command
ENTRYPOINT ["/app/homebox-export"]
CMD ["export"]

# Labels for container metadata
LABEL org.opencontainers.image.source="https://github.com/kusold/homebox-export" \
    org.opencontainers.image.title="homebox-export" \
    org.opencontainers.image.description="Export tool for Homebox" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.created="${DATE}" \
    org.opencontainers.image.revision="${COMMIT}" \
    maintainer="Mike Kusold <hello@mikekusold.com>"
# org.opencontainers.image.licenses="" \
