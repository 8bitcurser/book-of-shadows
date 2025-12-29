# Book of Shadows - Call of Cthulhu Character Sheet Manager
# Makefile for development, testing, and deployment

# Variables
APP_NAME := book-of-shadows
BINARY := $(APP_NAME)
GO := go
TEMPL := templ

# Build flags
CGO_ENABLED := 1
LDFLAGS := -s -w

# Fly.io settings
FLY_APP := book-of-shadows

.PHONY: all build run clean test lint fmt vet check deps dev watch deploy deploy-staging logs help

# Default target
all: build

# ============================================================================
# Development
# ============================================================================

## dev: Run the application in development mode
dev: generate
	$(GO) run .

## watch: Watch for templ changes and regenerate
watch:
	$(TEMPL) generate --watch

## generate: Generate templ templates
generate:
	$(TEMPL) generate

## run: Build and run the application
run: build
	./$(BINARY)

# ============================================================================
# Build
# ============================================================================

## build: Build the application binary
build: generate
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY) .

## build-linux: Build for Linux (used in Docker)
build-linux: generate
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY) .

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -f $(APP_NAME)-test
	find . -name "*_templ.go" -type f -delete

# ============================================================================
# Code Quality
# ============================================================================

## test: Run tests
test:
	$(GO) test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## lint: Run golangci-lint
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GO) fmt ./...
	$(TEMPL) fmt .

## vet: Run go vet
vet:
	$(GO) vet ./...

## check: Run all code quality checks
check: fmt vet lint test

# ============================================================================
# Dependencies
# ============================================================================

## deps: Download and tidy dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

## deps-update: Update all dependencies
deps-update:
	$(GO) get -u ./...
	$(GO) mod tidy

## tools: Install development tools
tools:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# ============================================================================
# Deployment (Fly.io)
# ============================================================================

## deploy: Deploy to Fly.io production
deploy: check
	@echo "Deploying to Fly.io..."
	fly deploy --app $(FLY_APP)

## deploy-now: Deploy to Fly.io without checks (use with caution)
deploy-now:
	@echo "Deploying to Fly.io (skipping checks)..."
	fly deploy --app $(FLY_APP)

## deploy-staging: Deploy to staging environment
deploy-staging:
	@echo "Deploying to staging..."
	fly deploy --app $(FLY_APP)-staging

## logs: View Fly.io logs
logs:
	fly logs --app $(FLY_APP)

## logs-follow: Follow Fly.io logs in real-time
logs-follow:
	fly logs --app $(FLY_APP) -f

## status: Check Fly.io app status
status:
	fly status --app $(FLY_APP)

## ssh: SSH into Fly.io machine
ssh:
	fly ssh console --app $(FLY_APP)

## scale: Show current scaling
scale:
	fly scale show --app $(FLY_APP)

## secrets: List Fly.io secrets
secrets:
	fly secrets list --app $(FLY_APP)

# ============================================================================
# Docker
# ============================================================================

## docker-build: Build Docker image locally
docker-build:
	docker build -t $(APP_NAME):latest .

## docker-run: Run Docker container locally
docker-run: docker-build
	docker run -p 8080:8080 -v $(PWD)/data:/data $(APP_NAME):latest

# ============================================================================
# Database
# ============================================================================

## db-backup: Backup the production database
db-backup:
	@mkdir -p backups
	fly ssh sftp get /data/exports.db backups/exports-$(shell date +%Y%m%d-%H%M%S).db --app $(FLY_APP)

## db-shell: Open SQLite shell on production
db-shell:
	fly ssh console --app $(FLY_APP) -C "sqlite3 /data/exports.db"

# ============================================================================
# Help
# ============================================================================

## help: Show this help message
help:
	@echo "Book of Shadows - Available Commands:"
	@echo ""
	@echo "Development:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | grep -E "^(dev|watch|generate|run):" | column -t -s ':'
	@echo ""
	@echo "Build:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | grep -E "^(build|build-linux|clean):" | column -t -s ':'
	@echo ""
	@echo "Code Quality:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | grep -E "^(test|lint|fmt|vet|check):" | column -t -s ':'
	@echo ""
	@echo "Deployment:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | grep -E "^(deploy|logs|status|ssh|scale):" | column -t -s ':'
	@echo ""
	@echo "Docker:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | grep -E "^docker-" | column -t -s ':'
	@echo ""
	@echo "Database:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | grep -E "^db-" | column -t -s ':'
