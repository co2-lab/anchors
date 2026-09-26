#!/usr/bin/env bash
# The unit suite as `anchors test` runs it: the JUnit report proves the scenarios, and the
# lcov report measures the lines. Both are declared in anchors.yaml (`tests:`), and a tool
# that is missing only drops its own report — the suite still runs.
set -euo pipefail

cd "$(dirname "$0")/.."
mkdir -p .anchors

tool() {
  if command -v "$1" >/dev/null 2>&1; then
    command -v "$1"
  elif [ -x "$HOME/go/bin/$1" ]; then
    echo "$HOME/go/bin/$1"
  fi
}
JUNIT_CMD="$(tool go-junit-report)"
LCOV_CMD="$(tool gcov2lcov)"

if [ -n "$JUNIT_CMD" ]; then
  go test -v -coverpkg=./... -coverprofile=.anchors/cover.out ./... | "$JUNIT_CMD" > .anchors/junit.xml
else
  go test -coverpkg=./... -coverprofile=.anchors/cover.out ./...
fi

if [ -n "$LCOV_CMD" ]; then
  "$LCOV_CMD" -infile .anchors/cover.out -outfile .anchors/lcov.info
fi
