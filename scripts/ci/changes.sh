#!/usr/bin/env bash
# Decides which CI areas a change needs, from the list of changed paths.
# ci.yml's `changes` job feeds it the paths dorny/paths-filter lists and gates
# every job on the outputs; see "What CI runs" in docs/TESTING.md.
#
# Rules are the ordered table below; the first pattern that matches a path
# wins. Patterns are bash [[ == ]] globs, where * also matches '/'. A path that
# matches no rule is "unclassified" and selects every area, so a new kind of
# file can never make CI skip work. Each rule sets flags:
#   b backend   f frontend   c compose only   w workflows (actionlint)
#   m Go modules (govulncheck)   p Postgres tests regardless of the Go closure
#   - ignored (no CI value)
# compose is derived: backend || frontend || c.
#
# Assumption: no *.md file under backend/ or web/ is embedded or imported. If
# one ever is, add an exception for it above the "*.md" rule.
#
# Usage (from anywhere; paths are repo-relative, one per line on stdin):
#   changes.sh classify      print backend=, frontend=, compose=, workflows=,
#                            gomod= lines
#   changes.sh go-closure    print postgres=: whether the paths touch the Go
#                            dependency closure of the Postgres test packages
#                            (scripts/ci/test-shards.json) and ./cmd/migrate
#   changes.sh check [--go]  run scripts/ci/changes-fixtures.txt; --go also runs
#                            the go-closure fixtures (needs Go and modules)
#
# Environment:
#   CI_EVENT         GitHub event name. The web/package.json version-only rule
#                    applies only to pull_request, whose checkout is the merge
#                    commit, so HEAD^1 is the base.
#   CHANGES_SUMMARY  if set, a Markdown table explaining each path is appended
#                    to this file (ci.yml passes $GITHUB_STEP_SUMMARY).
#
# Locally: git diff --name-only origin/main... | scripts/ci/changes.sh classify
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fixtures="$repo_root/scripts/ci/changes-fixtures.txt"
shards="$repo_root/scripts/ci/test-shards.json"

die() { echo "changes: $*" >&2; exit 1; }

# pattern<TAB>flags, in match order.
rules=$(cat <<'EOF'
.github/workflows/ci.yml	bfwmp
docs/openapi.yaml	f
*.md	-
docs/*	-
graphify-out/*	-
.graphifyignore	-
openspec/*	-
.claude/*	-
.agents/*	-
.codex/*	-
.idea/*	-
prototypes/*	-
deploy/*	-
LICENSE	-
.github/CODEOWNERS	-
.github/ISSUE_TEMPLATE/*	-
.github/dependabot.yml	-
.gitignore	-
cliff.toml	-
skills-lock.json	-
.air.toml	-
lefthook.yml	-
Makefile	-
CHANGELOG.md	-
release-please-config.json	-
.release-please-manifest.json	-
backend/go.mod	bmp
backend/go.sum	bmp
go.work	bmp
go.work.sum	bmp
backend/*	b
.golangci.yml	b
scripts/ci/*	bp
.github/actions/go-cache-*	bwp
web/*	f
.github/actions/bun-setup/*	fw
Dockerfile	c
.dockerignore	c
docker-compose*.yml	c
.env.example	c
scripts/compose-smoke.*	c
.github/workflows/*	w
.github/actions/*	w
.github/codeql/*	w
EOF
)

# True when web/package.json differs from the PR base only in "version".
# CHANGES_PKG_VERSION_ONLY (true/false) overrides git, for fixtures.
package_version_only() {
	if [[ -n ${CHANGES_PKG_VERSION_ONLY:-} ]]; then
		[[ $CHANGES_PKG_VERSION_ONLY == true ]]
		return
	fi
	[[ ${CI_EVENT:-} == pull_request ]] || return 1
	local base head
	base=$(git -C "$repo_root" show HEAD^1:web/package.json 2>/dev/null | jq -S 'del(.version)' 2>/dev/null) || return 1
	head=$(jq -S 'del(.version)' "$repo_root/web/package.json" 2>/dev/null) || return 1
	[[ -n $base && $base == "$head" ]]
}

# Prints "<rule><TAB><flags>" for one path; flags "u" = unclassified.
classify_path() {
	local path=$1 pat flags
	if [[ $path == web/package.json ]]; then
		if package_version_only; then
			printf 'web/package.json (version only)\t-\n'
		else
			printf 'web/package.json\tf\n'
		fi
		return
	fi
	while IFS=$'\t' read -r pat flags; do
		# shellcheck disable=SC2053 # $pat is a glob on purpose
		if [[ $path == $pat ]]; then
			printf '%s\t%s\n' "$pat" "$flags"
			return
		fi
	done <<<"$rules"
	printf 'unclassified\tu\n'
}

# Human-readable areas for a flag string.
describe_flags() {
	local flags=$1 out=()
	case $flags in
	-) echo "ignored"; return ;;
	u) echo "all (unclassified)"; return ;;
	esac
	[[ $flags == *b* ]] && out+=(backend)
	[[ $flags == *f* ]] && out+=(frontend)
	[[ $flags == *[bfc]* ]] && out+=(compose)
	[[ $flags == *w* ]] && out+=(workflows)
	[[ $flags == *m* ]] && out+=(gomod)
	[[ $flags == *p* ]] && out+=(postgres)
	local IFS=,
	echo "${out[*]}"
}

# Reads paths on stdin, prints the area outputs.
cmd_classify() {
	local path rule flags all=""
	local table="| Path | Rule | Selects |"$'\n'"|---|---|---|"$'\n'
	while IFS= read -r path; do
		[[ -n $path ]] || continue
		IFS=$'\t' read -r rule flags < <(classify_path "$path")
		all+=$flags
		table+="| \`$path\` | \`$rule\` | $(describe_flags "$flags") |"$'\n'
	done

	local backend=false frontend=false compose=false workflows=false gomod=false
	[[ $all == *[bu]* ]] && backend=true
	[[ $all == *[fu]* ]] && frontend=true
	[[ $backend == true || $frontend == true || $all == *c* ]] && compose=true
	[[ $all == *w* ]] && workflows=true
	[[ $all == *[mu]* ]] && gomod=true

	printf 'backend=%s\nfrontend=%s\ncompose=%s\nworkflows=%s\ngomod=%s\n' \
		"$backend" "$frontend" "$compose" "$workflows" "$gomod"

	if [[ -n ${CHANGES_SUMMARY:-} ]]; then
		{
			echo "### Change detection"
			echo
			echo "backend=$backend · frontend=$frontend · compose=$compose · workflows=$workflows · gomod=$gomod"
			echo
			echo "<details><summary>Why, per changed path</summary>"
			echo
			printf '%s' "$table"
			echo
			echo "</details>"
			echo
		} >>"$CHANGES_SUMMARY"
	fi
}

# Repo-relative Go package directories, one per line, tagged "pkg" for every
# in-module package and "closure" for the Postgres test packages' dependencies.
go_package_dirs() {
	local roots
	mapfile -t roots < <(jq -r '[.postgres[].packages[]] | unique | .[]' "$shards")
	cd "$repo_root/backend"
	go list -f 'pkg {{.Dir}}' ./...
	go list -deps -test -f '{{if not .Standard}}closure {{.Dir}}{{end}}' "${roots[@]}" ./cmd/migrate
}

# Reads paths on stdin, prints postgres=true|false. A path touches the closure
# when its nearest enclosing Go package is in it; that covers embedded files
# (migrations), testdata and deleted files. When in doubt (a forced rule, an
# unclassified path, no enclosing package, go list failing) it answers true.
cmd_go_closure() {
	local paths path rule flags dir listing
	mapfile -t paths
	local affected=false reason=""
	declare -A pkg=() closure=()

	for path in "${paths[@]}"; do
		[[ -n $path ]] || continue
		IFS=$'\t' read -r rule flags < <(classify_path "$path")
		if [[ $flags == *[pu]* ]]; then
			affected=true reason="\`$path\` ($rule)"
			break
		fi
	done

	if [[ $affected == false ]]; then
		if ! listing=$(go_package_dirs); then
			affected=true reason="go list failed"
		else
			local tag
			while read -r tag dir; do
				dir=${dir#"$repo_root"/}
				[[ -n $dir ]] || continue
				if [[ $tag == pkg ]]; then pkg[$dir]=1; else closure[$dir]=1; fi
			done <<<"$listing"
			for path in "${paths[@]}"; do
				[[ $path == backend/* ]] || continue
				dir=$(dirname "$path")
				while [[ $dir == backend/* && -z ${pkg[$dir]:-} ]]; do
					dir=$(dirname "$dir")
				done
				if [[ -z ${pkg[$dir]:-} ]]; then
					affected=true reason="\`$path\` (no enclosing Go package)"
					break
				fi
				if [[ -n ${closure[$dir]:-} ]]; then
					affected=true reason="\`$path\` (package ${dir#backend/})"
					break
				fi
			done
		fi
	fi

	echo "postgres=$affected"
	if [[ -n ${CHANGES_SUMMARY:-} ]]; then
		if [[ $affected == true ]]; then
			echo "Postgres tests: **run**, because of $reason." >>"$CHANGES_SUMMARY"
		else
			echo "Postgres tests: **skipped**, no changed path is in their Go dependency closure." >>"$CHANGES_SUMMARY"
		fi
	fi
}

# Fixture lines:
#   <path><TAB><areas>[<TAB>version-only]   areas: comma list of the true
#                                           classify outputs, or "-" for none
#   @postgres<TAB><path><TAB><true|false>   go-closure fixture (check --go)
cmd_check() {
	local with_go=false
	[[ ${1:-} == --go ]] && with_go=true
	local failed=0 n=0 line path want flag got
	while IFS= read -r line; do
		[[ -z $line || $line == \#* ]] && continue
		if [[ $line == @postgres* ]]; then
			[[ $with_go == true ]] || continue
			IFS=$'\t' read -r _ path want <<<"$line"
			got=$(printf '%s\n' "$path" | CHANGES_SUMMARY='' cmd_go_closure)
			got=${got#postgres=}
		else
			IFS=$'\t' read -r path want flag <<<"$line"
			local pkg_override=""
			[[ ${flag:-} == version-only ]] && pkg_override=true
			got=$(printf '%s\n' "$path" |
				CHANGES_SUMMARY='' CHANGES_PKG_VERSION_ONLY=${pkg_override:-false} cmd_classify |
				sed -n 's/=true$//p' | paste -sd, -)
			got=${got:--}
		fi
		n=$((n + 1))
		if [[ $got != "$want" ]]; then
			echo "changes: fixture mismatch: $line" >&2
			echo "changes:   got $got" >&2
			failed=1
		fi
	done <"$fixtures"
	((failed == 0)) || exit 1
	echo "changes: $n fixtures OK"
}

case ${1:-} in
classify) cmd_classify ;;
go-closure) cmd_go_closure ;;
check) shift; cmd_check "$@" ;;
*) die "usage: changes.sh classify | go-closure | check [--go]" ;;
esac
