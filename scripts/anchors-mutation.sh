#!/usr/bin/env bash
# The mutation run as `anchors mutation` runs it: gremlins over every package of the
# module, merged into the single report declared in anchors.yaml (`mutation:`).
#
# One gremlins run per package, because a run over a subpackage writes each `file_name`
# relative to THAT package: `queue.go`, not `internal/queue/queue.go`. Ingested as is, a
# bare name matches every file of the map with that basename, and `cmd/anchors/flow/queue.go`
# would receive the score of `internal/queue/queue.go`. Each report is prefixed with its
# package directory before the merge.
#
# gremlins also DESCENDS into subdirectories: `./cmd/anchors` mutated every package under
# it, which each has its own run here. `-E /` keeps each run to the files of its own
# package (the regexp is matched against paths relative to the target, and only a file of
# a subdirectory has a `/` in it).
#
# `--timeout-coefficient 10`: gremlins times each mutant against the coverage run, which
# is fast when the build cache is warm. A mutant that needs a rebuild after the cache was
# trimmed then exceeds it, and a TIMED OUT mutant counts as killed — measured on
# internal/testlist: 32 of 32 timed out (a false 100%), and 30 killed / 1 survived with
# the larger coefficient.
#
# ANCHORS_MUTATION_TIMEOUT_COEFFICIENT overrides the 10: when mutation-score reports a file
# measured under load, measure it once with a large coefficient (100) to learn how long a
# mutant really takes, as the gate says.
#
# Usage: scripts/anchors-mutation.sh [package-dir | file.go ...]   (default: every package)
set -uo pipefail

cd "$(dirname "$0")/.."
mkdir -p .anchors
out=.anchors/mutation.json
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# Each entry is "package-dir" or "package-dir|file.go". A `.go` argument (relative, or
# absolute under the repository, as `run_changed: {{files}}` passes it) mutates that file
# alone: gremlins takes a package, so its sibling files are excluded one by one.
root="$(pwd)"
if [ "$#" -gt 0 ]; then
  pkgs=()
  for a in "$@"; do
    a="${a#"$root"/}"
    case "$a" in
      *.go) pkgs+=("$(dirname "$a")|$(basename "$a")") ;;
      *) pkgs+=("$a") ;;
    esac
  done
else
  mod="$(go list -m)"
  pkgs=()
  while IFS= read -r p; do
    pkgs+=("${p#"$mod"/}")
  done < <(go list ./... | grep -v "^$mod\$")
fi

coefficient="${ANCHORS_MUTATION_TIMEOUT_COEFFICIENT:-10}"
i=0
for entry in "${pkgs[@]}"; do
  i=$((i + 1))
  pkg="${entry%%|*}"
  excludes=(-E /)
  if [ "$entry" != "$pkg" ]; then
    only="${entry#*|}"
    for sibling in "$pkg"/*.go; do
      name="$(basename "$sibling")"
      [ "$name" = "$only" ] && continue
      excludes+=(-E "^${name//./\\.}\$")
    done
  fi
  echo "[$i/${#pkgs[@]}] $entry"
  if ! gremlins unleash "./$pkg" "${excludes[@]}" --timeout-coefficient "$coefficient" --output "$tmp/raw.json" >"$tmp/log" 2>&1; then
    echo "  gremlins failed on $pkg:" >&2
    tail -5 "$tmp/log" >&2
    continue
  fi
  grep -E 'Test efficacy' "$tmp/log" | sed 's/^/  /'
  jq --arg dir "$pkg" '.files |= map(.file_name = ($dir + "/" + .file_name))' \
    "$tmp/raw.json" >"$tmp/pkg-$i.json"
done

shopt -s nullglob
reports=("$tmp"/pkg-*.json)
if [ "${#reports[@]}" -eq 0 ]; then
  echo "no package produced a report" >&2
  exit 1
fi
jq -s '{go_module: .[0].go_module, files: (map(.files) | add)}' "${reports[@]}" >"$out"
echo "merged ${#reports[@]} package report(s) into $out"
