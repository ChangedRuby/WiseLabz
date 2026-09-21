#!/usr/bin/env bash
set -Eeuo pipefail

project_name="${COMPOSE_PROJECT_NAME:-wiselabz-smoke}"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
env_file="$(mktemp)"
port="${COMPOSE_SMOKE_PORT:-18080}"
export COMPOSE_SMOKE_ENV_FILE="$env_file"
export COMPOSE_SMOKE_PORT="$port"
compose=(docker compose --project-name "$project_name" --env-file "$env_file" -f "$repo_root/docker-compose.yml" -f "$repo_root/scripts/compose-smoke.yml")

cleanup() {
	status=$?
	if (( status != 0 )); then
		"${compose[@]}" ps || true
		"${compose[@]}" logs --no-color || true
	fi
	"${compose[@]}" down --volumes --remove-orphans || true
	rm -f "$env_file"
	exit "$status"
}
trap cleanup EXIT

cat >"$env_file" <<EOF
WISELABZ_AUTH_SECRET=compose-smoke-auth-secret-0123456789abcdef
WISELABZ_ENCRYPTION_KEY=$(openssl rand -base64 32)
WISELABZ_SERVER_ORIGIN=http://127.0.0.1:$port
WISELABZ_ADMIN_PASSWORD=compose-smoke-admin-password
POSTGRES_USER=wiselabz
POSTGRES_PASSWORD=compose-smoke-postgres-password
POSTGRES_DB=wiselabz
EOF

"${compose[@]}" config -q
"${compose[@]}" up --build --wait --wait-timeout 120

base_url="http://127.0.0.1:$port/api"
root_url="${base_url%/api}"
curl --fail --silent --show-error "$root_url/readyz" | jq -e '.ready == true' >/dev/null
access_token="$(curl --fail --silent --show-error \
	-H 'Content-Type: application/json' \
	--data '{"username":"admin","password":"compose-smoke-admin-password"}' \
	"$base_url/auth/login" | jq -er '.accessToken')"
connector_id="$(curl --fail --silent --show-error \
	-H "Authorization: Bearer $access_token" \
	-H 'Content-Type: application/json' \
	--data '{"name":"Compose smoke Docker","category":"containers_paas","type":"docker","url":"tcp://docker.invalid:2375","config":{"host":"tcp://docker.invalid:2375"}}' \
	"$base_url/connectors" | jq -er '.id')"
curl --fail --silent --show-error \
	-H "Authorization: Bearer $access_token" \
	"$base_url/connectors" | jq -e --arg id "$connector_id" 'any(.[]; .id == $id and .name == "Compose smoke Docker")' >/dev/null
