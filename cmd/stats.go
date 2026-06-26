package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Your full focus time",
	Long:  "See how much foucsed time you had during the day",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		csvFile := "pomodoro.csv"
		if os.Getenv("csvfile") != "" {
			csvFile = os.Getenv("csvfile")
		}
		filename := filepath.Join(homeDir, "Library", "Application Support", "pomodoro", csvFile)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		var total time.Duration
		for _, r := range records {
			d, err := time.ParseDuration(r[1])
			if err != nil {
				log.Fatal(err)
			}
			total += d
		}
		fmt.Printf("%s\n", total.String())
	},
}
