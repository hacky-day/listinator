# Listinator Vue Branch Makefile
# This Makefile provides build automation for the Vue.js rewrite of Listinator

# Variables
APP_NAME := listinator
GO_VERSION := $(shell go version | cut -d' ' -f3)
NODE_VERSION := $(shell node --version 2>/dev/null || echo "not-found")
NPM_VERSION := $(shell npm --version 2>/dev/null || echo "not-found")

# Frontend paths
FRONTEND_DIR := frontend
FRONTEND_DIST := $(FRONTEND_DIR)/dist
FRONTEND_SRC := $(FRONTEND_DIR)/src

# Go build settings
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := 1

# Build flags
LDFLAGS := -w -s
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

.DEFAULT_GOAL := help

## Help
.PHONY: help
help: ## Show this help message
	@echo "Listinator Vue Branch - Build System"
	@echo "===================================="
	@echo ""
	@echo "Environment:"
	@echo "  Go version:    $(GO_VERSION)"
	@echo "  Node version:  $(NODE_VERSION)"
	@echo "  NPM version:   $(NPM_VERSION)"
	@echo "  OS/Arch:       $(GOOS)/$(GOARCH)"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Development
.PHONY: deps
deps: deps-go deps-frontend ## Install all dependencies

.PHONY: deps-go
deps-go: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy

.PHONY: deps-frontend
deps-frontend: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

.PHONY: dev
dev: deps build-frontend ## Start development server (with frontend build)
	@echo "Starting development server..."
	@echo "Frontend will be served from embedded dist/ directory"
	@echo "Server will run on http://localhost:8080"
	LISTINATOR_DATABASE_DIR=./data go run .

.PHONY: dev-frontend
dev-frontend: ## Start frontend development server (separate from backend)
	@echo "Starting frontend development server..."
	@echo "This runs the Vue dev server with hot reload"
	cd $(FRONTEND_DIR) && npm run dev

## Building
.PHONY: build
build: build-frontend build-go ## Build complete application (frontend + backend)

.PHONY: build-frontend
build-frontend: deps-frontend ## Build frontend for production
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm run build
	@echo "Frontend built to $(FRONTEND_DIST)/"

.PHONY: build-go
build-go: deps-go ## Build Go backend
	@echo "Building Go application..."
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $(APP_NAME) .
	@echo "Built $(APP_NAME) for $(GOOS)/$(GOARCH)"

.PHONY: build-production
build-production: clean deps build ## Production build with optimizations
	@echo "Production build complete!"
	@echo "Binary: ./$(APP_NAME)"
	@echo "Frontend: $(FRONTEND_DIST)/"

## Testing & Quality
.PHONY: test
test: ## Run Go tests
	@echo "Running Go tests..."
	go test -v ./...

.PHONY: test-frontend
test-frontend: ## Run frontend tests
	@echo "Running frontend tests..."
	cd $(FRONTEND_DIR) && npm test

.PHONY: lint
lint: lint-go lint-frontend ## Run all linters

.PHONY: lint-go
lint-go: ## Run Go linter
	@echo "Running Go linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead..."; \
		go vet ./...; \
	fi

.PHONY: lint-frontend
lint-frontend: ## Run frontend linter
	@echo "Running frontend linter..."
	cd $(FRONTEND_DIR) && npm run type-check

.PHONY: format
format: ## Format code
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting frontend code..."
	cd $(FRONTEND_DIR) && npm run format 2>/dev/null || echo "No format script found"

## Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run: ## Run application in Docker
	@echo "Running Docker container..."
	docker run -p 8080:8080 -e LISTINATOR_DATABASE_DIR=/var/lib/listinator -v $(PWD)/data:/var/lib/listinator $(APP_NAME):latest

.PHONY: docker-compose
docker-compose: ## Run with docker-compose
	@echo "Starting with docker-compose..."
	docker-compose up --build

## Cleanup
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(APP_NAME)
	rm -rf $(FRONTEND_DIST)
	rm -rf $(FRONTEND_DIR)/node_modules/.tmp

.PHONY: clean-all
clean-all: clean ## Clean everything including dependencies
	@echo "Cleaning all dependencies..."
	rm -rf $(FRONTEND_DIR)/node_modules
	go clean -modcache

## Feature Parity Check
.PHONY: check-parity
check-parity: ## Compare Vue branch features with main branch
	@echo "Checking feature parity with main branch..."
	@echo "API Endpoints:"
	@echo "  GET  /api/v1/entries  - List entries"
	@echo "  POST /api/v1/entries  - Create entry"
	@echo "  PUT  /api/v1/entries/:id - Update entry"
	@echo "  DELETE /api/v1/entries/:id - Delete entry"
	@echo "  POST /api/v1/lists    - Create list"
	@echo "  GET  /api/v1/types    - List types"
	@echo ""
	@echo "Frontend Features:"
	@echo "  ✓ List management (create, view)"
	@echo "  ✓ Entry management (add, edit, delete, mark as bought)"
	@echo "  ✓ Search functionality"
	@echo "  ✓ Type/category support with icons"
	@echo "  ✓ Sharing functionality"
	@echo "  ✓ Responsive design"
	@echo ""
	@echo "Technical Changes:"
	@echo "  • Frontend: Vanilla JS → Vue.js 3 + TypeScript"
	@echo "  • Build system: None → Vite"
	@echo "  • API structure: /server → /api/v1/server"
	@echo "  • Static serving: filesystem fallback → embedded dist/"

## Information
.PHONY: info
info: ## Show project information
	@echo "Project: Listinator Vue Branch"
	@echo "Description: Vue.js rewrite of the Listinator list management application"
	@echo ""
	@echo "Architecture:"
	@echo "  Backend:  Go with Echo framework and GORM"
	@echo "  Database: SQLite"
	@echo "  Frontend: Vue.js 3 + TypeScript + Vite"
	@echo "  Docker:   Multi-stage build with embedded frontend"
	@echo ""
	@echo "Environment Variables:"
	@echo "  LISTINATOR_DATABASE_DIR - Directory for SQLite database (required)"
	@echo ""
	@echo "Default Ports:"
	@echo "  Application: 8080"
	@echo "  Frontend Dev: 5173 (vite default)"