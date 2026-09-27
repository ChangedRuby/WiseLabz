#!/usr/bin/env bash
# Runs one CI shard of a backend Go test suite. Shards are defined once in
# scripts/ci/test-shards.json:
#
#   { "<suite>": { "<shard>": { "packages": [...],
#                               "tests": [...] | "pattern": "..." | "rest": true,
#                               "exclude": [...] } } }
#
# - "tests": the shard runs only these top-level tests (-run '^(A|B)$'). Use it
#   when a few slow tests dominate a package.
# - "pattern": the shard runs top-level tests matching this regex (e.g. a name
#   range like '^Test[A-C]'). Use it when a package has many similar-sized tests.
# - "rest": the shard runs every test in its packages NOT selected by a sibling
#   shard with the same packages (-skip). New tests land here automatically, so
#   nothing is silently dropped when tests are added or renamed.
# - "exclude": packages removed from "packages" (they belong to another shard).
#
# Usage (from anywhere):
#   test-shards.sh matrix <suite>                     JSON array of shard names
#   test-shards.sh run <suite> <shard> [go test flags...] [-- test binary args...]
#   test-shards.sh check                              validate the definition file
#   test-shards.sh timings <suite> [go test flags...] per-test durations, slowest first
#
# Re-balancing: run `timings` (with the same flags CI uses, e.g. -race), then
# greedily assign the slowest tests to the currently lightest shard and update
# the JSON. CI timings are roughly 1.7x local ones.
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
defs="$repo_root/scripts/ci/test-shards.json"
cd "$repo_root/backend"

die() { echo "test-shards: $*" >&2; exit 1; }

shard_json() {
	jq -e --arg s "$1" --arg h "$2" '.[$s][$h]' "$defs" >/dev/null ||
		die "unknown shard $1/$2"
	jq -c --arg s "$1" --arg h "$2" '.[$s][$h]' "$defs"
}

# Resolved import paths for a shard: packages minus exclude.
shard_packages() {
	local shard=$1 excl pkgs p
	mapfile -t pkgs < <(jq -r '.packages[]' <<<"$shard" | xargs go list)
	excl="$(jq -r '.exclude // [] | .[]' <<<"$shard" | xargs -r go list)"
	for p in "${pkgs[@]}"; do
		grep -qxF "$p" <<<"$excl" || echo "$p"
	done
}

# jq: the top-level test regex a non-rest shard selects ("" = no filter).
shard_regex='def shard_regex: if has("tests") then .tests | map("^" + . + "$") | join("|") else .pattern // "" end;'

cmd_matrix() {
	jq -c --arg s "$1" '.[$s] | keys_unsorted' "$defs"
}

cmd_run() {
	local suite=$1 name=$2 shard filter=() pkgs regex
	shift 2
	shard="$(shard_json "$suite" "$name")"
	if jq -e '.rest == true' <<<"$shard" >/dev/null; then
		regex="$(jq -r --arg s "$suite" --argjson me "$shard" "$shard_regex"'
			[.[$s][] | select(.packages == $me.packages and .rest != true) | shard_regex] | join("|")' "$defs")"
		[[ -n "$regex" ]] && filter=(-skip "$regex")
	else
		regex="$(jq -r "$shard_regex"' shard_regex' <<<"$shard")"
		[[ -n "$regex" ]] && filter=(-run "$regex")
	fi

	local flags=() binargs=()
	while (($#)); do
		if [[ $1 == -- ]]; then shift; binargs=(-args "$@"); break; fi
		flags+=("$1"); shift
	done
	mapfile -t pkgs < <(shard_packages "$shard")
	((${#pkgs[@]})) || die "shard $suite/$name resolves to no packages"

	set -x
	go test "${flags[@]}" "${filter[@]}" "${pkgs[@]}" "${binargs[@]}"
}

cmd_check() {
	local status=0 suite
	for suite in $(jq -r 'keys[]' "$defs"); do
		# Every listed test must exist, and appear in exactly one shard.
		local dupes
		dupes="$(jq -r --arg s "$suite" '[.[$s][] | .tests // [] | .[]] | group_by(.) | map(select(length > 1) | .[0]) | .[]' "$defs")"
		if [[ -n "$dupes" ]]; then
			echo "$suite: tests listed in more than one shard: $dupes" >&2
			status=1
		fi
		local groups
		groups="$(jq -c --arg s "$suite" '[.[$s][] | select(has("tests") or has("pattern"))] | group_by(.packages)[] | {packages: .[0].packages, tests: [.[].tests // [] | .[]]}' "$defs")"
		while IFS= read -r g; do
			[[ -n "$g" ]] || continue
			if ! jq -e --arg s "$suite" --argjson g "$g" \
				'[.[$s][] | select(.rest == true and .packages == $g.packages)] | length == 1' "$defs" >/dev/null; then
				echo "$suite: shards over $(jq -c .packages <<<"$g") need exactly one \"rest\" shard" >&2
				status=1
			fi
			local have
			have="$(jq -r '.packages[]' <<<"$g" | xargs go test -list '.*' | grep -E '^(Test|Example|Fuzz)' | sort -u)"
			while IFS= read -r t; do
				if ! grep -qxF "$t" <<<"$have"; then
					echo "$suite: $t is listed in a shard but does not exist" >&2
					status=1
				fi
			done < <(jq -r '.tests[]' <<<"$g")
		done <<<"$groups"
		# Every excluded package must be run by a sibling shard.
		local covered
		covered="$(jq -r --arg s "$suite" '[.[$s][] | select(has("exclude") | not) | .packages[]] | .[]' "$defs" | xargs -r go list)"
		while IFS= read -r p; do
			[[ -n "$p" ]] || continue
			if ! grep -qxF "$p" <<<"$covered"; then
				echo "$suite: excluded package $p is not run by any shard" >&2
				status=1
			fi
		done < <(jq -r --arg s "$suite" '[.[$s][] | .exclude // [] | .[]] | .[]' "$defs" | xargs -r go list)
	done
	((status == 0)) && echo "test-shards: definitions OK"
	return "$status"
}

cmd_timings() {
	local suite=$1
	shift
	mapfile -t pkgs < <(jq -r --arg s "$suite" '[.[$s][] | .packages[]] | unique | .[]' "$defs")
	go test -json -count=1 "$@" "${pkgs[@]}" |
		jq -r 'select(.Action == "pass" or .Action == "fail") | select(.Test != null and (.Test | contains("/") | not)) | "\(.Elapsed)\t\(.Package | split("/") | last)\t\(.Test)"' |
		sort -rn
}

case "${1:-}" in
matrix) (($# == 2)) || die "usage: matrix <suite>"; cmd_matrix "$2" ;;
run) (($# >= 3)) || die "usage: run <suite> <shard> [flags...]"; shift; cmd_run "$@" ;;
check) cmd_check ;;
timings) (($# >= 2)) || die "usage: timings <suite> [flags...]"; shift; cmd_timings "$@" ;;
*) die "usage: test-shards.sh matrix|run|check|timings ..." ;;
esac
