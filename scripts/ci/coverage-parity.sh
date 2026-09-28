#!/usr/bin/env bash
# Compares per-package statement coverage of two Go text coverage profiles
# (go test -coverprofile, or `go tool covdata textfmt`) and fails when any
# package differs by more than a tolerance in percentage points. Used to check
# that a cheaper coverage strategy keeps -coverpkg=./... cross-package
# attribution (e.g. router tests covering handler packages), not just the
# total.
#
# Usage: coverage-parity.sh <baseline.out> <candidate.out> [tolerance-pp=0.1]
#
# NOISY_PACKAGES (regex) get NOISY_TOLERANCE (pp) instead: their coverage
# varies between identical runs because some branches only execute under
# particular goroutine interleavings. internal/ws moves ~0.6pp run to run
# (ws.go client write-pump select branches), so it is allowed 1pp by default;
# a real attribution loss is far larger than that.
set -Eeuo pipefail

(($# >= 2)) || { echo "usage: coverage-parity.sh <baseline.out> <candidate.out> [tolerance-pp]" >&2; exit 2; }

# "<package>\t<covered>\t<total>" per package. A block can appear several times
# (one per test binary under -coverpkg); count it once, covered if any hit it.
per_package() {
	awk 'NR > 1 {
		split($1, loc, ":"); file = loc[1]
		pkg = file; sub(/\/[^\/]*$/, "", pkg)
		if (!($1 in stmts)) { stmts[$1] = $2; owner[$1] = pkg }
		if ($3 > 0) hit[$1] = 1
	}
	END {
		for (b in stmts) { total[owner[b]] += stmts[b]; if (b in hit) cov[owner[b]] += stmts[b] }
		for (p in total) printf "%s\t%d\t%d\n", p, cov[p], total[p]
	}' "$1" | sort
}

join -t$'\t' -a1 -a2 -e0 -o 0,1.2,1.3,2.2,2.3 \
	<(per_package "$1") <(per_package "$2") |
	awk -F'\t' -v tol="${3:-0.1}" \
		-v noisy="${NOISY_PACKAGES-/internal/ws\$}" -v noisytol="${NOISY_TOLERANCE:-1}" '
	function pct(c, t) { return t ? 100 * c / t : 0 }
	{
		b = pct($2, $3); c = pct($4, $5); d = c - b
		bc += $2; bt += $3; cc += $4; ct += $5
		lim = (noisy != "" && $1 ~ noisy) ? noisytol : tol
		flag = (d > lim || d < -lim) ? "  <-- differs" : ""
		if (flag != "") bad++
		if (flag != "" || ENVIRON["VERBOSE"] != "") printf "%-70s %6.1f%% -> %6.1f%% (%+.1fpp)%s\n", $1, b, c, d, flag
	}
	END {
		printf "TOTAL %.1f%% -> %.1f%%; %d package(s) differ by more than %spp\n", pct(bc, bt), pct(cc, ct), bad, tol
		exit bad > 0
	}'
