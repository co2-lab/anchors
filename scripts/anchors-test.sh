#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
mkdir -p .anchors

JUNIT_CMD="go-junit-report"
if ! command -v go-junit-report >/dev/null 2>&1; then
  if [ -x "$HOME/go/bin/go-junit-report" ]; then
    JUNIT_CMD="$HOME/go/bin/go-junit-report"
  fi
fi

if command -v "$JUNIT_CMD" >/dev/null 2>&1 || [ -x "$JUNIT_CMD" ]; then
  go test -v ./... | "$JUNIT_CMD" > .anchors/junit.xml
else
  go test ./...
fi
