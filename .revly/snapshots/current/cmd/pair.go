package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nareshkarthigeyan/revly/internals/cache"
	"github.com/nareshkarthigeyan/revly/internals/gitutils"
	"github.com/nareshkarthigeyan/revly/internals/llm"
	"github.com/spf13/cobra"
)

var pairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Start a pair programming session with Revly.",
	Long:  `Watches your code for changes and provides live AI feedback on your current diffs.`,
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetInt("interval")
		log.Printf("Starting pair programming mode with %d second interval...", interval)

		// Create .revly directory if it doesn't exist
		if _, err := os.Stat(".revly"); os.IsNotExist(err) {
			os.Mkdir(".revly", 0755)
		}

		// Set up logging
		logFile, err := os.OpenFile(".revly/pair.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)

		// Create initial snapshot
		initialSnapshotPath := ".revly/snapshots/initial"
		if err := gitutils.CreateSnapshot(initialSnapshotPath); err != nil {
			log.Fatalf("Failed to create initial snapshot: %v", err)
		}

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Create current snapshot
			currentSnapshotPath := ".revly/snapshots/current"
			if err := gitutils.CreateSnapshot(currentSnapshotPath); err != nil {
				log.Printf("Failed to create current snapshot: %v", err)
				continue
			}

			diff, err := gitutils.DiffSnapshots(initialSnapshotPath, currentSnapshotPath)
			if err != nil {
				fmt.Printf("Error getting diff: %v\n", err)
				continue
			}

			if len(diff) == 0 {
				// fmt.Println("No changes detected.")
				continue
			}

			key := cache.Key(diff)
			if _, err := cache.Load(key); err == nil {
				// fmt.Println("No new changes detected.")
				continue
			}

			comment, err := llm.GetPairProgrammingComment(string(diff))
			if err != nil {
				fmt.Printf("Error getting comment: %v\n", err)
				continue
			}

			if comment != "" {
				fmt.Printf("Suggestion: %s\n", comment)
				log.Printf("Suggestion: %s", comment)
				cache.Save(key, []byte(comment))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pairCmd)
	pairCmd.Flags().IntP("interval", "i", 15, "Interval in seconds to check for changes.")
}