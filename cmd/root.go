package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/57blocks/guardian-side/internal"
)

const (
	linuxPath  = "geth/guardian"
	darwinPath = "Library/Story/geth/guardian"
)

// Global variable to hold the output directory.
var outputDir string

// rootCmd is the root command for the Guardian side.
var rootCmd = &cobra.Command{
	Use:   "guardian-side",
	Short: "A tool to periodically download bloom filter files.",
	Run: func(cmd *cobra.Command, args []string) {
		startBloomDownloadTask()
	},
}

// Initializes the command-line flags and bind them to Viper configurations.
func init() {
	// Register the output directory flag and bind it with Viper.
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", getDefaultPath(), "Directory to store the bloom filter file")
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output-dir"))
}

// Execute is the main entry point to start the Cobra CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execution failed, err: %v", err)
	}
}

// startBloomDownloadTask initializes a periodic task to download the bloom filter file once per day.
func startBloomDownloadTask() {
	downloadAndRetry()

	for {
		// Calculate the time to next midnight.
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		sleepDuration := nextMidnight.Sub(now)

		// Sleep until next midnight to start the next download.
		time.Sleep(sleepDuration)

		// Retry and download the file again after the sleep period.
		downloadAndRetry()
	}
}

// downloadAndRetry downloads the bloom filter file with a retry mechanism.
func downloadAndRetry() {
	err := retry.Do(
		func() error {
			// Attempt to download and store bloom filter
			if err := internal.DownloadAndSaveBloomFilter(outputDir); err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			return nil
		},
		retry.Delay(3*time.Second),
		retry.Attempts(6),
	)
	if err != nil {
		log.Printf("Failed to download bloom filter after retries: %v", err)
	} else {
		log.Printf("Successfully downloaded bloom filter to %s", outputDir)
	}
}

// getDefaultPath determines the default file path based on the operating system.
func getDefaultPath() string {
	userHomeDir, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "linux":
		return filepath.Join(userHomeDir, linuxPath)
	case "darwin":
		return filepath.Join(userHomeDir, darwinPath)
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
		return ""
	}
}
