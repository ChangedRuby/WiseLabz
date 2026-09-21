# syntax=docker/dockerfile:1

# =============================================================================
# Stage 1/3: frontend — build the Vite/React/TS SPA -> web/dist/
# =============================================================================
FROM oven/bun:1-alpine AS frontend
WORKDIR /repo/web

# Install deps first, isolated from source changes, so `bun install` is cached.
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile

# `bun run build` runs orval codegen via the npm "prebuild" lifecycle hook
# before tsc/vite, and orval.config.ts points at ../docs/openapi.yaml — so
# the spec must be present at that exact relative path before building.
COPY web/ ./
COPY docs/openapi.yaml /repo/docs/openapi.yaml
RUN bun run build

# =============================================================================
# Stage 2/3: backend — build the static Go binary
# =============================================================================
FROM golang:1.27-alpine AS backend

# git is required for `go build`'s automatic VCS stamping (debug.ReadBuildInfo,
# included in the operator diagnostics bundle) — the alpine golang image ships without it.
# ca-certificates is required for TLS to the Go module proxy during download.
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy module files and resolve deps first, so this layer is cached unless
# go.mod/go.sum change.
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download

# The Go module lives in backend/, not the repo root, but VCS stamping only
# activates when the main package, its module, and the working directory all
# resolve to the same git repository. So .git must be copied as a sibling of
# backend/ here (mirroring the real repo layout), not just the backend
# subtree on its own.
COPY .git ./.git
COPY backend/ ./backend/

# Overlay the real frontend build over the placeholder committed at
# backend/internal/web/dist/ (see backend/internal/web/embed.go) so the
# `//go:embed all:dist` picks up the actual SPA instead of the placeholder.
COPY --from=frontend /repo/web/dist/. ./backend/internal/web/dist/

WORKDIR /src/backend

# CGO_ENABLED=0 is safe: the backend is 100% pure-Go (modernc.org/sqlite for
# SQLite, jackc/pgx/v5/stdlib for Postgres) — no cgo/libc dependency, so the
# resulting binary is statically linked and runs on the cgo-less distroless
# base in stage 3.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/wiselabz ./cmd/server

# Pre-create the default SQLite data directory owned by the distroless
# "nonroot" UID/GID (65532:65532). distroless/static has no shell to mkdir
# or chown, so this has to happen in the (shell-having) builder stage and be
# copied in; Docker also uses an existing image-layer directory's ownership
# as the template when a volume is later mounted at the same path, which is
# what makes `docker run -v data:/data` work out of the box for the default
# file:/data/wiselabz.db DSN.
RUN mkdir -p /data && chown -R 65532:65532 /data

# =============================================================================
# Stage 3/3: final — minimal distroless runtime image
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot AS final

COPY --from=backend /out/wiselabz /wiselabz
COPY --from=backend --chown=65532:65532 /data /data

USER nonroot:nonroot

EXPOSE 8080
VOLUME ["/data"]

# Serve the embedded SPA build (server.embed); override at runtime if needed.
ENV WISELABZ_SERVER_EMBED=true

# distroless/static has no shell (no curl/wget), so HEALTHCHECK relies on the
# binary's own --healthcheck mode: it GETs its own /readyz and exits 0 only
# when the server is ready.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD ["/wiselabz", "--healthcheck"]

ENTRYPOINT ["/wiselabz"]
