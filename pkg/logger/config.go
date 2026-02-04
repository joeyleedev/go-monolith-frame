package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds logging configuration
type Config struct {
	// Basic configuration
	Level string `mapstructure:"level"` // debug, info, warn, error

	// Output configuration
	Filename    string `mapstructure:"filename"`    // Log file path
	Console     bool   `mapstructure:"console"`     // Enable console output
	EnableColor bool   `mapstructure:"enable_color"` // Enable color output

	// File rotation configuration (lumberjack)
	MaxSize    int  `mapstructure:"max_size"`    // Single log file max size before rotation (MB)
	MaxAge     int  `mapstructure:"max_age"`     // Time to live (days)
	MaxBackups int  `mapstructure:"max_backups"` // Max log file count
	Compress   bool `mapstructure:"compress"`    // Compress rotated files
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Level:       "info",
		Console:     true,
		EnableColor: true,
	}
}

// Validate validates configuration
func (c *Config) Validate() error {
	// Validate log level
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[c.Level] {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	// Validate filename
	if c.Filename != "" {
		dir := filepath.Dir(c.Filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create log directory: %v", err)
		}
	}

	return nil
}
