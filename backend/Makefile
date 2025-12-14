# Backend Makefile
# Champions League Simulation API

.PHONY: all build run test coverage lint lint-fix format format-check deps verify docker-build docker-run swagger swagger-install

## run: Run the application
run: build
	@echo "Running..."
	@go run ./cmd/server/main.go

## test: Run tests
test:
	@echo "Running tests..."
	@go test -race -v ./...

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## lint-fix: Run linter with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix ./...

## format: Format code with golangci-lint (uses gofumpt)
format:
	@echo "Formatting code..."
	@golangci-lint fmt ./...
	@echo "Formatted"


## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded"

## swagger: Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	@swag init -g ./cmd/server/main.go -o ./internal/handlers/docs --parseDependency --parseInternal
	@echo "Swagger docs generated at ./internal/handlers/docs"

## swagger-install: Install swag CLI tool
swagger-install:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@v1.16.6
	@echo "Swag installed"


## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t champions-league-backend .

## docker-run: Run Docker container
docker-run:
	@docker run -p 8080:8080 champions-league-backend

review: swagger lint format test verify
	@echo "Review completed successfully!"