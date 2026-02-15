#!/usr/bin/env bash
set -euo pipefail

# Quick wrapper around assetviz that supports either file input or stdin.
# It keeps core tool logic untouched and only makes invocation easier.

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  cat <<'USAGE'
Usage:
  ./scripts/quickviz.sh subdomains.txt
  cat subdomains.txt | ./scripts/quickviz.sh

Behavior:
  - If an argument is provided, runs: go run . -f <file>
  - If stdin is piped, runs: go run .
USAGE
  exit 0
fi

if [[ $# -gt 0 ]]; then
  go run . -f "$1"
else
  if [ -t 0 ]; then
    echo "No input received. Pass a file or pipe data via stdin." >&2
    exit 1
  fi
  go run .
fi
