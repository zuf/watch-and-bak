#!/bin/bash

set -e

# Sleep and skip all changes in files for that period (seconds)
# after every backup operation (to not create backups too often).
SLEEP_PERIOD=60
BASE_BAK_DIR='./bak/'

while true;
do
	FILE=$(inotifywait -e modify $@ | awk '{print $1}')	
	BASENAME=$(basename "${FILE}")
	BAK_PREFIX=$(date +%Y%m%d_%H%M%S_)
	FILE_BAK_DIR="${BASE_BAK_DIR}/$(date +%Y-%m-%d)/"	
	FULL_BAK_PATH="${FILE_BAK_DIR}/${BAK_PREFIX}${BASENAME}.gz"
	
	mkdir -p "${FILE_BAK_DIR}"
	cat "$FILE" | gzip -9 > "${FULL_BAK_PATH}"
	
	echo "File ${BASENAME} saved to ${FULL_BAK_PATH}" > /dev/stderr
	
	sleep $SLEEP_PERIOD
done