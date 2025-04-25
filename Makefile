# Simple Makefile for a Go project

include .env

# Build the application
all: build

build:
	@echo "Building..."
	@templ generate
	@npx @tailwindcss/cli -i ./web/assets/css/input.css -o ./web/assets/css/tailwind.css --minify
	@go build -o piggy-planner main.go 

build-release:
	@echo "Building..."
	@templ generate
	@npx @tailwindcss/cli -i ./web/assets/css/input.css -o ./web/assets/css/tailwind.css --minify
	@go build -ldflags="-s -w" -o piggy-planner main.go 

build-release-raspberry:
	@echo "Building..."
	@templ generate
	@npx @tailwindcss/cli -i ./web/assets/css/input.css -o ./web/assets/css/tailwind.css --minify
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o piggy-planner-raspberry main.go 

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
