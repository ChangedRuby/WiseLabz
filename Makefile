.PHONY: dev-backend dev-frontend build test lint migrate-up migrate-down hooks setup dev fmt help

dev-backend:
	@which air > /dev/null 2>&1 && air -c .air.toml || go run ./backend/cmd/server

dev-frontend:
	npm --prefix web run dev

build:
	npm --prefix web run build
	rm -rf backend/internal/web/dist
	mkdir -p backend/internal/web/dist
	cp -r web/dist/. backend/internal/web/dist/
	(cd backend && go build -o ../bin/wiselabz ./cmd/server)

test:
	go test -short ./backend/...

lint:
	golangci-lint run ./backend/...
	npm --prefix web run lint

migrate-up:
	(cd backend && go run ./cmd/migrate up)

migrate-down:
	(cd backend && go run ./cmd/migrate down)

hooks:
	lefthook install

# Install every dependency and wire up the git hooks. Run once after cloning.
setup:
	npm --prefix web install
	cd backend && go mod download
	$(MAKE) hooks

# Run the backend and the frontend dev servers side by side.
dev:
	$(MAKE) -j2 dev-backend dev-frontend

# Format the Go backend and the web front-end.
fmt:
	gofmt -w backend
	npm --prefix web run format

# List the available make targets.
help:
	@awk '/^[a-zA-Z0-9_-]+:/ {name=$$1; sub(/:.*/, "", name); printf "  %-16s\n", name}' $(MAKEFILE_LIST)
