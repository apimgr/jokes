.PHONY: build build-all build-host local dev release docker test check clean help

PROJECTNAME := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)(\.git)?$$|\1|' || basename "$(PWD)")
PROJECTORG := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)/[^/]+(\.git)?$$|\1|' || basename "$(dir $(PWD))")
VERSION ?= $(shell cat release.txt 2>/dev/null || echo "0.1.0")
COMMIT_ID := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo unknown)
BUILD_EPOCH := $(shell date -u +%s)
BUILD_DATE := $(shell date -u -d @$(BUILD_EPOCH) +%Y-%m-%dT%H:%M:%SZ)
OFFICIAL_SITE ?= https://jokes.apimgr.us
BINDIR := binaries
RELDIR := releases
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 freebsd/amd64 freebsd/arm64

GO_CACHE ?= $(HOME)/go/pkg/mod
GO_BUILD ?= $(HOME)/.cache/go-build/$(PROJECTNAME)
TMPDIR ?= /tmp
DOCKER_MEM ?= 4g
DOCKER_CPUS ?= 2
GO_DOCKER := docker run --rm \
	--name $(PROJECTNAME)-$$(tr -dc 'a-z0-9' </dev/urandom | head -c8) \
	--memory=$(DOCKER_MEM) --cpus=$(DOCKER_CPUS) \
	-v "$(PWD):/app" \
	-v "$(GO_CACHE):/usr/local/share/go/pkg/mod" \
	-v "$(GO_BUILD):/usr/local/share/go/cache" \
	-v "$(TMPDIR)/$(PROJECTORG):/tmp/$(PROJECTORG)" \
	-w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false \
	casjaysdev/go:latest
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.CommitID=$(COMMIT_ID) -X main.BuildDate=$(BUILD_DATE) -X main.BuildEpoch=$(BUILD_EPOCH) -X main.OfficialSite=$(OFFICIAL_SITE)

## build: Build both binaries for the required release platforms and host
build: build-all build-host

## build-all: Cross-build server and CLI for all eight supported platforms
build-all:
	@mkdir -p $(GO_CACHE) $(GO_BUILD) $(BINDIR) $(RELDIR)
	@$(GO_DOCKER) sh -c 'set -eu; for platform in $(PLATFORMS); do os=$${platform%/*}; arch=$${platform#*/}; ext=""; [ "$$os" != windows ] || ext=.exe; GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" -o "$(BINDIR)/$(PROJECTNAME)-$$os-$$arch$$ext" ./src; GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" -o "$(BINDIR)/$(PROJECTNAME)-cli-$$os-$$arch$$ext" ./src/client; done'

## build-host: Build server and CLI for the current host platform
build-host:
	@mkdir -p $(GO_CACHE) $(GO_BUILD) $(BINDIR)
	@$(GO_DOCKER) sh -c 'go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(PROJECTNAME) ./src && go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(PROJECTNAME)-cli ./src/client'

## local: Build server and CLI for production testing
local: build-host

## dev: Build server and CLI into a project-scoped temporary directory
dev:
	@mkdir -p $(GO_CACHE) $(GO_BUILD) "$(TMPDIR)/$(PROJECTORG)"
	@$(GO_DOCKER) sh -c 'set -eu; out=$$(mktemp -d "/tmp/$(PROJECTORG)/$(PROJECTNAME)-XXXXXX"); go build -buildvcs=false -o "$$out/$(PROJECTNAME)" ./src; go build -buildvcs=false -o "$$out/$(PROJECTNAME)-cli" ./src/client; printf "%s\n" "$$out"'

## release: Create versioned release archives locally without pushing or tagging
release: build
	@mkdir -p $(RELDIR)
	@tar -C $(BINDIR) -czf $(RELDIR)/$(PROJECTNAME)-$(VERSION)-binaries.tar.gz .
	@tar --exclude=.git --exclude=binaries --exclude=releases -czf $(RELDIR)/$(PROJECTNAME)-$(VERSION)-source.tar.gz .

## docker: Build the local container image without publishing it
docker:
	@docker buildx build --platform linux/amd64 --load --tag ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION) --file docker/Dockerfile .

## test: Run Go vet, vulnerability checks, and the test suite in Docker
test:
	@mkdir -p $(GO_CACHE) $(GO_BUILD)
	@$(GO_DOCKER) sh -c 'mkdir -p /tmp/$(PROJECTORG) && covdir=$$(mktemp -d /tmp/$(PROJECTORG)/$(PROJECTNAME)-XXXXXX) && trap "rm -rf $$covdir" EXIT && go vet ./... && govulncheck ./... && go test -coverprofile=$$covdir/coverage.out ./... && coverage=$$(go tool cover -func=$$covdir/coverage.out | awk "/^total:/ {gsub(/%/, \"\", \$$3); print \$$3}") && echo "Coverage: $$coverage%" && awk -v coverage="$$coverage" "BEGIN { if (coverage+0 < 60) exit 1 }"'

## check: Run static analysis and formatting checks in Docker
check:
	@mkdir -p $(GO_CACHE) $(GO_BUILD)
	@$(GO_DOCKER) sh -c 'set -e; go vet ./...; GOTOOLCHAIN=go1.26.4 staticcheck ./...; unformatted=$$(gofmt -l ./src); if [ -n "$$unformatted" ]; then printf "Unformatted Go files:\n%s\n" "$$unformatted"; exit 1; fi'

## clean: Remove generated binaries, release archives, and Go cache for this project
clean:
	@rm -rf -- $(BINDIR) $(RELDIR) "$(GO_BUILD)"

## help: Show available targets
help:
	@printf '%s\n' 'Jokes Makefile' '' 'Targets:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
	@printf '\nVersion: %s\n' '$(VERSION)'
