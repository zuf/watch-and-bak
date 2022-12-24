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

