#!/usr/bin/env bash

# https://peter.sh/experiments/chromium-command-line-switches/

CHROME_DATA_DIR=$(mktemp -d)
trap "rm -rf ${CHROME_DATA_DIR}" SIGINT SIGTERM EXIT

PREVURL=$(cat /home/pi/urlfile.txt)
OLIVEURL="https://oliveai.com"
URL=${PREVURL:-$OLIVEURL}

/usr/bin/chromium-browser --remote-debugging-port=9222 --user-data-dir="${CHROME_DATA_DIR}" --disable-infobars --kiosk "${URL}"
