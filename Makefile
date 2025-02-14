GOLANG          := golang:1.24
ALPINE          := alpine:3.21
HAMERKOP_APP    := hamerkop
BASE_IMAGE_NAME := localhost/charlieroth
VERSION         := 0.1.0
HAMERKOP_IMAGE  := $(BASE_IMAGE_NAME)/$(HAMERKOP_APP):$(VERSION)

# ==============================================================================
# Install dependencies

dev-docker:
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
	wait;

# ==============================================================================
# Build containers 

build: hamerkop

hamerkop:
	docker build \
		-f zarf/docker/dockerfile.hamerkop \
		-t $(HAMERKOP_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Docker Compose 

compose-up:
	cd .zarf/compose && docker compose -f docker-compose.yaml -p compose up -d

compose-down:
	cd .zarf/compose && docker compose -f docker-compose.yaml down

compose-logs:
	cd .zarf/compose && docker compose -f docker-compose.yaml logs

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Local Development

run:
	go run cmd/hamerkop/hamerkop.go

clean:
	rm -f hamerkop

nip11:
	curl -H "Accept: application/nostr+json" http://localhost:8080 | jq