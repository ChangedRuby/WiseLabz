#!/usr/bin/env bash
# Prints the -coverpkg list for backend coverage: every package (./...) so
# router-level tests in internal/api attribute coverage to the handler, store
# and domain packages they exercise, minus test-only helper packages whose
# statements would otherwise count toward the total.
#
# Usage (from backend/): go test -coverpkg="$(../scripts/ci/coverpkg.sh)" ./...
set -Eeuo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/../../backend"
go list ./... | grep -vE '/internal/(api/apitest|store/storetest|connector/connectortest)$' | paste -sd, -
