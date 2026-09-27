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
# Usage: scripts/anchors-mutation.sh [package-dir ...]   (default: every package)
set -uo pipefail

cd "$(dirname "$0")/.."
mkdir -p .anchors
out=.anchors/mutation.json
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

if [ "$#" -gt 0 ]; then
  pkgs=("$@")
else
  mod="$(go list -m)"
  pkgs=()
  while IFS= read -r p; do
    pkgs+=("${p#"$mod"/}")
  done < <(go list ./... | grep -v "^$mod\$")
fi

i=0
for pkg in "${pkgs[@]}"; do
  i=$((i + 1))
  echo "[$i/${#pkgs[@]}] $pkg"
  if ! gremlins unleash "./$pkg" -E / --output "$tmp/raw.json" >"$tmp/log" 2>&1; then
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
