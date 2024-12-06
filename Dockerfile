FROM alpine:3.21
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

# Copy binary from goreleaser
COPY --chown=appuser:appgroup homebox-export /app/homebox-export

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
