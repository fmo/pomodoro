package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetCsvFilename() string {
	csvFile := "pomodoro.csv"

	if os.Getenv("csvfile") != "" {
		csvFile = os.Getenv("csvfile")
	}

	return csvFile
}

func GetCsvBackup(csvfile string) string {
	s := strings.Split(csvfile, ".")
	return fmt.Sprintf("%s_bup.csv", s[0])
}

func GetFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "pomodoro")
}

func CreateBackup(file *os.File) error {
	src, err := os.Open(file.Name())
	if err != nil {
		return err
	}

	fp := GetFilePath()
	csvFile := GetCsvBackup(GetCsvFilename())

	dsc, err := os.Create(filepath.Join(fp, csvFile))
	if err != nil {
		return err
	}

	io.Copy(dsc, src)

	return nil
}

func Save(limit, count int) error {
	fp := GetFilePath()

	err := os.MkdirAll(fp, 0o700)
	if err != nil {
		return err
	}

	fullFileName := filepath.Join(fp, GetCsvFilename())

	var f *os.File

	f, err = os.OpenFile(fullFileName, os.O_APPEND|os.O_WRONLY, 0o700)
	if err != nil {
		f, err = os.Create(fullFileName)
		if err != nil {
			return err
		}
	}

	if cfg != nil {
		if bup, ok := cfg.Get("app.backup-active").(bool); ok {
			if bup {
				CreateBackup(f)
			}
		}
	}

	w := csv.NewWriter(f)

	// completed is total duration in seconds - count down time, ie total 30 seconds - 10 count down (left)
	completed := time.Duration(limit-count) * time.Second

	left := strconv.Itoa(count)

	// current time, completed time, count down
	err = w.Write([]string{time.Now().Format(time.RFC3339), completed.String(), fmt.Sprintf("%ss", left)})
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}
