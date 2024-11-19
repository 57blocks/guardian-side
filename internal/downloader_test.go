package internal

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestDownloadAndSaveBloomFilter(t *testing.T) {
	type args struct {
		outputDir string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		mockResp  string
		httpError bool
	}{
		{
			name:      "Successful download and save",
			args:      args{outputDir: os.TempDir()},
			wantErr:   false,
			mockResp:  "bloom_filter_data",
			httpError: false,
		},
		{
			name:      "Network error in downloading",
			args:      args{outputDir: os.TempDir()},
			wantErr:   true,
			httpError: true,
		},
		{
			name:     "Invalid output directory",
			args:     args{outputDir: "/invalid-path"},
			wantErr:  true,
			mockResp: "bloom_filter_data",
		},
	}

	// Initialize HTTPMock globally for this test.
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock based on test case configuration.
			if tt.httpError {
				// Simulate HTTP error
				httpmock.RegisterResponder("GET", bloomFilterFileURL,
					httpmock.NewErrorResponder(errors.New("error")))
			} else {
				// Simulate successful response from remote server.
				httpmock.RegisterResponder("GET", bloomFilterFileURL,
					httpmock.NewStringResponder(http.StatusOK, tt.mockResp))
			}

			// Run the function under test.
			err := DownloadAndSaveBloomFilter(tt.args.outputDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadAndSaveBloomFilter() error = %v, wantErr %v", err, tt.wantErr)
			}

			// If the test should succeed, check if the file was written correctly.
			if !tt.wantErr {
				filePath := filepath.Join(tt.args.outputDir, bloomFilterFilename)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file does not exist: %v", filePath)
				}
				defer os.Remove(filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read downloaded file: %v", err)
				}

				if string(content) != tt.mockResp {
					t.Errorf("Expected file content %s but got %s", tt.mockResp, string(content))
				}
			}
		})
	}
}
