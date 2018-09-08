#!/usr/bin/env bash

# https://peter.sh/experiments/chromium-command-line-switches/

CHROME_DATA_DIR=$(mktemp -d)
trap "rm -rf ${CHROME_DATA_DIR}" SIGINT SIGTERM EXIT
chromium --remote-debugging-port=9222 --user-data-dir="${CHROME_DATA_DIR}"

