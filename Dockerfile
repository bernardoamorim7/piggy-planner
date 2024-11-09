# Stage 1: Download and unzip the Go binary
FROM alpine:latest AS downloader

# Install wget and unzip
RUN apk add --no-cache wget unzip

# Set the working directory
WORKDIR /app

# Download the zip file from GitHub releases
ADD https://github.com/bernardoamorim7/piggy-planner/releases/download/v0.1.1/piggy-planner_0.1.1_linux_arm64.zip /app/piggy-planner.zip

# Unzip the binary
RUN unzip piggy-planner.zip

# Make the binary executable
RUN chmod +x piggy-planner

# Stage 2: Create a minimal image with the Go binary
FROM scratch

# Set the working directory
WORKDIR /app

# Add necessary certificates for HTTPS connections
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the Go binary from the downloader stage
COPY --from=downloader /app/piggy-planner /app/piggy-planner

# Copy the local file into the container
# Uncomment the line below if you want to copy the database file into the container
# COPY ./piggy_planner.db /app/piggy_planner.db

# Set the entrypoint to the Go binary
ENTRYPOINT ["/app/piggy-planner"]