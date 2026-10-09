#!/usr/bin/env bash
# Fleet sweep: run the securitymd provider across fleet repositories.
#
# Preview by default (detect only). --fix generates a SECURITY.md where none
# exists; BuildFlow's provider never overwrites an existing policy, so the
# sweep is safe to re-run.
#
# Usage:
#   scripts/fleet-securitymd-sweep.sh repo [repo...]
#   scripts/fleet-securitymd-sweep.sh --repos fleet-repos.txt [--fix]
#
# Each repo argument is a path to a checked-out fleet repository. The repos
# file lists one path per line (# comments and blank lines allowed).

set -euo pipefail

usage() {
	echo "usage: $0 [--fix] (--repos FILE | repo [repo...])" >&2
	exit 2
}

fix=false
repos_file=""
repos=()

while [[ $# -gt 0 ]]; do
	case "$1" in
	--fix)
		fix=true
		shift
		;;
	--repos)
		[[ $# -ge 2 ]] || usage
		repos_file="$2"
		shift 2
		;;
	-*)
		usage
		;;
	*)
		repos+=("$1")
		shift
		;;
	esac
done

if [[ -n "$repos_file" ]]; then
	[[ ${#repos[@]} -eq 0 ]] || usage
	while IFS= read -r line; do
		line="${line%%#*}"
		line="$(echo "$line" | tr -d '[:space:]')"
		[[ -n "$line" ]] && repos+=("$line")
	done <"$repos_file"
fi

[[ ${#repos[@]} -ge 1 ]] || usage

if ! command -v buildflow >/dev/null 2>&1; then
	echo "error: buildflow not found on PATH (build it in the BuildFlow repository first)" >&2
	exit 2
fi

mode="detect (preview)"
$fix && mode="fix (generate where missing)"
echo "securitymd fleet sweep — $mode, ${#repos[@]} repo(s)"
echo

repaired=0
clean=0
failed=0

for repo in "${repos[@]}"; do
	if [[ ! -d "$repo" ]]; then
		echo "SKIP  $repo (not a directory)"
		failed=$((failed + 1))
		continue
	fi

	echo "── $repo"
	if $fix; then
		buildflow -s securitymd --fix || {
			echo "FAIL  $repo"
			failed=$((failed + 1))
			continue
		}
	else
		buildflow -s securitymd || {
			echo "FAIL  $repo"
			failed=$((failed + 1))
			continue
		}
	fi

	if [[ -f "$repo/SECURITY.md" || -f "$repo/.github/SECURITY.md" || -f "$repo/docs/SECURITY.md" ]]; then
		clean=$((clean + 1))
		echo "OK    $repo"
	else
		repaired=$((repaired + 1))
		echo "REPAIR-NEEDED  $repo (no policy in any candidate location)"
	fi
done

echo
echo "summary: $clean with policy, $repaired without, $failed failed/skipped"
[[ $failed -eq 0 ]]
