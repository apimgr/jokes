#!/usr/bin/env bash
set -euo pipefail

if (($# == 0)); then
    set -- /usr/local/bin/jokes --port "${PORT:-80}"
fi

exec "$@"
