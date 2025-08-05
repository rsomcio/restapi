.PHONY: test test-unit test-integration build run clean

# Run all tests
test:
	go test -v ./...

# Run only unit tests (without database)
test-unit:
	go test -v ./models ./handlers -short

# Run integration tests (requires DATABASE_URL)
test-integration:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL environment variable required for integration tests"; \
		exit 1; \
	fi
	go test -v ./database

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the application
build:
	go build -o bin/restapi main.go

# Run the application (requires DATABASE_URL)
run:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL environment variable required"; \
		exit 1; \
	fi
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run all checks (format, lint, test)
check: fmt lint test