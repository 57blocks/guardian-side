package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config represents the configuration structure.
type Config struct {
	S3Bucket string `mapstructure:"s3_bucket"` // S3 bucket name
	S3Key    string `mapstructure:"s3_key"`    // S3 object key
	Region   string `mapstructure:"region"`    // AWS region
	FilePath string `mapstructure:"file_path"` // File path to read from
}

// readConfig initializes and retrieves the configuration using Viper.
func readConfig() (*Config, error) {
	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Look for the config in the current directory
	viper.AutomaticEnv()          // Automatically bind environment variables
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}
