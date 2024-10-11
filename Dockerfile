# Use an official MariaDB image as the base
FROM mariadb:lts

# Set environment variables for MariaDB
ENV MARIADB_DATABASE=piggy_planner
ENV MARIADB_USER=${DB_USERNAME}
ENV MARIADB_PASSWORD=${DB_PASSWORD}
ENV MARIADB_ROOT_PASSWORD=${DB_ROOT_PASSWORD}

# Install necessary tools
RUN apt-get update && apt-get install -y curl jq

# Set the working directory
WORKDIR /app

# Download the latest release binary from GitHub
RUN LATEST_RELEASE_URL=$(curl -s https://api.github.com/repos/yourusername/yourrepository/releases/latest | jq -r '.assets[] | select(.name | contains("linux_amd64")) | .browser_download_url') && \
    curl -L -o main $LATEST_RELEASE_URL && \
    chmod +x piggy_planner

# Expose the application port
EXPOSE ${PORT}

# Start both MariaDB and the application
CMD ["sh", "-c", "mysqld_safe & /app/piggy_planner"]