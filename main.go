package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Conf struct {
	PollPeriod time.Duration
	// BackupDir     string
	BakFilePrefix string
	GzipLevel     int
}

type WatchFile struct {
	// TODO
	// sha1 sum
	//
}

func Usage() {
	fmt.Printf("Usage:\n  %s [OPTIONS] /path/to/file [/path/to/another/file] ...\n\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func InitConf() *Conf {
	conf := new(Conf)

	flag.DurationVar(&conf.PollPeriod, "n", 60*time.Second,
		`Specify polling interval. A duration string is a sequence of decimal
numbers, each with optional fraction and a unit suffix, such as "1m30s"
or "-1.5h". Valid time units are "s", "m", "h"`)
	// flag.StringVar(&conf.BackupDir, "d", "./bak", "Specify backup directory")
	flag.IntVar(&conf.GzipLevel, "z", gzip.BestCompression, fmt.Sprintf("Specify gzip compression level. Values from %d to %d.", gzip.BestSpeed, gzip.BestCompression))

	fimeFormatHelp := `Format explanation:	
  Year: "2006" "06"
  Month: "Jan" "January" "01" "1"
  Day of the week: "Mon" "Monday"
  Day of the month: "2" "_2" "02"
  Day of the year: "__2" "002"
  Hour: "15" "3" "03" (PM or AM)
  Minute: "4" "04"
  Second: "5" "05"
  AM/PM mark: "PM"

`
	flag.StringVar(&conf.BakFilePrefix, "p", "20060102_150405_", "Specify backup file prefix\n"+fimeFormatHelp)

	flag.Usage = Usage

	return conf
}

func Sha1Sum(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

func FileStatOrDie(filePath string) os.FileInfo {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return fileInfo
}

func MakeBackupDirForFile(filePath string) string {
	now := time.Now()

	// baseBakDir, err := filepath.Abs(conf.BackupDir)F
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO Ability to change "bak" dir in conf
	baseBakDir, err := filepath.Abs(filepath.Join(filepath.Dir(absFilePath), "bak"))
	if err != nil {
		log.Fatal(err)
	}

	subDir := now.Format("2006-01-02")
	bakDir := filepath.Join(baseBakDir, subDir)

	err = os.MkdirAll(bakDir, 0o755)
	if err != nil {
		log.Fatalf("Can't create directory \"%s\": %s\n", bakDir, err)
	}

	return bakDir
}

func GzCopyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	zw, err := gzip.NewWriterLevel(destination, conf.GzipLevel)
	if err != nil {
		return err
	}

	zw.ModTime = time.Now()

	_, err = io.Copy(zw, source)

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func BackupFile(filePath string) {
	bakDir := MakeBackupDirForFile(filePath)
	now := time.Now()
	bakFilePath := filepath.Join(bakDir, now.Format(conf.BakFilePrefix)+filepath.Base(filePath)+".gz")
	err := GzCopyFile(filePath, bakFilePath)
	if err != nil {
		log.Fatalf("Can't copy file \"%s\" to %s: %s\n", filePath, bakFilePath, err)
	}
	// TODO print bakFilePath path from basedir istead of absolute path?
	// TODO print human readable file size
	log.Printf("Backup of file %s was created at %s\n", filePath, bakFilePath)
}

var conf *Conf

func main() {
	conf = InitConf()
	flag.Parse()

	if flag.NArg() <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	// TODO Use Watchfile type in list instead this maps
	prevHashSums := make(map[int][]byte)
	prevFileStats := make(map[int]os.FileInfo)

	fmt.Printf("Watch and backup files:\n")
	for index, file := range flag.Args() {
		fmt.Printf("  %s\n", file)

		prevHashSums[index] = Sha1Sum(file)
		prevFileStats[index] = FileStatOrDie(file)
	}

	for {
		for index, filePath := range flag.Args() {

			prevFileStat := prevFileStats[index]
			fileStat := FileStatOrDie(filePath)

			if prevFileStat.ModTime() != fileStat.ModTime() {
				prevFileStat = fileStat
				prevFileStats[index] = prevFileStat

				prevHashSum := prevHashSums[index]
				hashSum := Sha1Sum(filePath)

				if !bytes.Equal(prevHashSum, hashSum) {
					prevHashSum = hashSum
					prevHashSums[index] = prevHashSum

					BackupFile(filePath)
				}
			}
		}
		time.Sleep(conf.PollPeriod)
	}
}
