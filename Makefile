.SHELL := /usr/bin/env bash
 
BINARY_NAME := dotfiles
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE)

.PHONY: build 
build:
	go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: run 
 run: 
	go run ./cmd/$(BINARY_NAME)/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf bin/$(BINARY_NAME)

.PHONY: patch minor major release-patch release-minor release-major release-semver
patch: release-patch

minor: release-minor

major: release-major

release-patch:
	@$(MAKE) release-semver BUMP=patch

release-minor:
	@$(MAKE) release-semver BUMP=minor

release-major:
	@$(MAKE) release-semver BUMP=major

# Create and publish the next semantic version tag.
# A GitHub Actions workflow is responsible for creating the GitHub Release from the pushed tag.
# Usage:
#   make release-patch
#   make release-minor
#   make release-major
release-semver:
	@set -euo pipefail; \
	if [[ -n "$$(git status --porcelain)" ]]; then \
		echo "Working tree is dirty. Commit or stash changes before creating a release."; \
		exit 1; \
	fi; \
	latest_tag="$$(git tag -l 'v*.*.*' --sort=-version:refname | head -n1)"; \
	if [[ -z "$$latest_tag" ]]; then \
		major=0; minor=0; patch=0; \
	else \
		version="$${latest_tag#v}"; \
		IFS='.' read -r major minor patch <<< "$$version"; \
	fi; \
	case "$(BUMP)" in \
		patch) patch=$$((patch + 1)); ;; \
		minor) minor=$$((minor + 1)); patch=0; ;; \
		major) major=$$((major + 1)); minor=0; patch=0; ;; \
		*) echo "Unsupported BUMP='$(BUMP)'. Use patch, minor, or major."; exit 1; ;; \
	esac; \
	next_tag="v$${major}.$${minor}.$${patch}"; \
	echo "Latest tag: $${latest_tag:-<none>}"; \
	echo "Creating release: $$next_tag"; \
	if git rev-parse -q --verify "refs/tags/$$next_tag" >/dev/null; then \
		echo "Tag $$next_tag already exists."; \
		exit 1; \
	fi; \
	git tag -a "$$next_tag" -m "$$next_tag"; \
	git push origin "$$next_tag"; \
	echo "Tag created and pushed: $$next_tag"