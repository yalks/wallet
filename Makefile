.PHONY: all build test clean fmt lint vet coverage deps help release-patch release-minor release-major quick-release

# Default target
all: fmt vet test build

# Build the module
build:
	@echo "Building..."
	@go build -v ./...

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./... || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean -cache
	@go clean -testcache

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	@go mod verify

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Check for security vulnerabilities
security:
	@echo "Checking for vulnerabilities..."
	@govulncheck ./... || echo "Install govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@godoc -http=:6060 & echo "Documentation server started at http://localhost:6060/pkg/github.com/yalks/wallet/"

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install golang.org/x/tools/cmd/godoc@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Release management
release-patch:
	@echo "Creating patch release..."
	@./release.sh patch

release-minor:
	@echo "Creating minor release..."
	@./release.sh minor

release-major:
	@echo "Creating major release..."
	@./release.sh major

quick-release:
	@echo "Quick patch release..."
	@./quick-release.sh

# Show help
help:
	@echo "Available targets:"
	@echo "  make all        - Format, vet, test, and build (default)"
	@echo "  make build      - Build the module"
	@echo "  make test       - Run tests"
	@echo "  make coverage   - Run tests with coverage report"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linter"
	@echo "  make vet        - Run go vet"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download and tidy dependencies"
	@echo "  make verify     - Verify dependencies"
	@echo "  make bench      - Run benchmarks"
	@echo "  make security   - Check for vulnerabilities"
	@echo "  make docs       - Start documentation server"
	@echo "  make check      - Run all checks (fmt, vet, lint, test)"
	@echo "  make install-tools - Install development tools"
	@echo ""
	@echo "Release targets:"
	@echo "  make quick-release  - Quick patch release (recommended)"
	@echo "  make release-patch  - Create patch release (x.x.X)"
	@echo "  make release-minor  - Create minor release (x.X.0)"
	@echo "  make release-major  - Create major release (X.0.0)"
	@echo ""
	@echo "  make help       - Show this help message"