.PHONY: build release docker test clean help

# Project variables
PROJECT_NAME := jokes
ORG := apimgr
VERSION := $(shell cat release.txt 2>/dev/null || echo "1.0.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build variables
BUILD_DIR := ./binaries
RELEASE_DIR := ./releases
BINARY_NAME := $(PROJECT_NAME)
DOCKER_IMAGE := ghcr.io/$(ORG)/$(PROJECT_NAME)

# Go build flags
LDFLAGS := -ldflags "-s -w -X main.VERSION=$(VERSION) -X main.BUILD_TIME=$(BUILD_TIME) -X main.GIT_COMMIT=$(GIT_COMMIT)"
CGO_ENABLED := 0

# Build targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 freebsd/amd64 freebsd/arm64

## build: Build all binaries for all platforms
build: clean
	@echo "🔨 Building $(PROJECT_NAME) v$(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR) $(RELEASE_DIR)
	@$(MAKE) build-all
	@$(MAKE) build-host
	@echo "✅ Build complete!"

## build-all: Build for all platforms
build-all:
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BUILD_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then \
			OUTPUT_NAME=$$OUTPUT_NAME.exe; \
		fi; \
		echo "📦 Building for $$GOOS/$$GOARCH..."; \
		CGO_ENABLED=$(CGO_ENABLED) GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $$OUTPUT_NAME . || exit 1; \
		if echo $$OUTPUT_NAME | grep -q "musl"; then \
			strip $$OUTPUT_NAME 2>/dev/null || true; \
		fi; \
	done

## build-host: Build for host platform
build-host:
	@echo "🏠 Building for host platform..."
	@CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Host binary: $(BUILD_DIR)/$(BINARY_NAME)"

## release: Create GitHub release with all binaries
release: build
	@echo "🚀 Creating release v$(VERSION)..."
	@mkdir -p $(RELEASE_DIR)
	@cp -r $(BUILD_DIR)/* $(RELEASE_DIR)/

	# Create source archive (no VCS files)
	@tar --exclude='.git' --exclude='node_modules' --exclude='binaries' --exclude='releases' \
		-czf $(RELEASE_DIR)/$(PROJECT_NAME)-$(VERSION)-source.tar.gz .

	# Check if tag exists and delete if it does
	@if git rev-parse v$(VERSION) >/dev/null 2>&1; then \
		echo "🗑️  Deleting existing tag v$(VERSION)..."; \
		git tag -d v$(VERSION); \
		git push --delete origin v$(VERSION) 2>/dev/null || true; \
	fi

	# Create new tag
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)

	# Create GitHub release
	@if command -v gh >/dev/null 2>&1; then \
		echo "📤 Creating GitHub release..."; \
		gh release create v$(VERSION) $(RELEASE_DIR)/* \
			--title "$(PROJECT_NAME) v$(VERSION)" \
			--notes "Release v$(VERSION) - Built on $(BUILD_TIME)"; \
	else \
		echo "⚠️  gh CLI not found. Please install it to create GitHub releases."; \
	fi

	# Increment version
	@$(MAKE) version-bump
	@echo "✅ Release v$(VERSION) created!"

## docker: Build and push Docker image
docker:
	@echo "🐳 Building Docker image..."
	@docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--push .
	@echo "✅ Docker image pushed: $(DOCKER_IMAGE):$(VERSION)"

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests complete! Coverage report: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(RELEASE_DIR)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "✅ Clean complete!"

## version-bump: Increment patch version
version-bump:
	@echo "📈 Incrementing version..."
	@CURRENT_VERSION=$$(cat release.txt); \
	MAJOR=$$(echo $$CURRENT_VERSION | cut -d'.' -f1); \
	MINOR=$$(echo $$CURRENT_VERSION | cut -d'.' -f2); \
	PATCH=$$(echo $$CURRENT_VERSION | cut -d'.' -f3); \
	NEW_PATCH=$$((PATCH + 1)); \
	NEW_VERSION="$$MAJOR.$$MINOR.$$NEW_PATCH"; \
	echo $$NEW_VERSION > release.txt; \
	echo "📌 Version bumped: $$CURRENT_VERSION → $$NEW_VERSION"

## help: Show this help message
help:
	@echo "$(PROJECT_NAME) Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Version: $(VERSION)"
