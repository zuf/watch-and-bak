# watch-and-bak

Very simple backup tool.
Watches for given file and make backups (.gz) in separate dir after that file was changed.

## Install

```
$ go install github.com/zuf/watch-and-bak@latest
```


## Usage

```
$ watch-and-bak --help

Usage:
  watch-and-bak [OPTIONS] /path/to/file [/path/to/another/file] ...

  -d string
    	Specify backup directory (default "./bak")
  -n duration
    	Specify polling interval. A duration string is a sequence of decimal
    	numbers, each with optional fraction and a unit suffix, such as "1m30s"
    	or "-1.5h". Valid time units are "s", "m", "h" (default 1m0s)
  -p string
    	Specify backup file prefix
    	Format explanation:	
    	  Year: "2006" "06"
    	  Month: "Jan" "January" "01" "1"
    	  Day of the week: "Mon" "Monday"
    	  Day of the month: "2" "_2" "02"
    	  Day of the year: "__2" "002"
    	  Hour: "15" "3" "03" (PM or AM)
    	  Minute: "4" "04"
    	  Second: "5" "05"
    	  AM/PM mark: "PM"
    	
    	 (default "20060102_150405_")
  -z int
    	Specify gzip compression level. Values from 1 to 9. (default 9)

```


## Alternative (on linux)

Do not want use weird software from internet?

Use plain old bash like that (you need to install `inotify-tools`):


```bash
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
```

Usage:


```
$ ./watch-and-bak.sh /path/to/your/files...
```

This script justan example, adapt for your needs.

[watch-and-bak.sh](watch-and-bak.sh)
