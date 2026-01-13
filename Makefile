.PHONY: build run test swagger migrate clean

# Build the application
build:
	go build -o bin/openauth cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -race ./...

# Run tests for specific package
test-package:
	@read -p "Enter package path: " pkg; \
	go test -v ./$$pkg

# Generate Swagger documentation
swagger:
	@echo "Installing/updating swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	@$$(go env GOPATH)/bin/swag init -g cmd/server/main.go -o docs/swagger || \
	 (echo "Trying with swag from PATH..." && swag init -g cmd/server/main.go -o docs/swagger)
	@echo "✅ Swagger documentation generated in docs/swagger/"

# Run database migrations
migrate:
	go run cmd/server/main.go migrate

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/swagger/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run all checks
check: fmt lint test
