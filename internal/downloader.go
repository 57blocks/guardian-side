package internal

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	bloomFilterFilename = "bloom_filter.gob"
	bloomFilterFileURL  = "https://remote-server1.com/bloomfilter"
)

// DownloadAndSaveBloomFilter downloads the remote bloom filter file and
// saves it to the specified location.
func DownloadAndSaveBloomFilter(outputDir string) error {
	// Check or create output directory
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	resp, err := http.Get(bloomFilterFileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write download content to file
	filePath := filepath.Join(outputDir, bloomFilterFilename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
