package config

import (
	"encoding/json"
	"fmt" // Added fmt import
	"io/ioutil"
	"os"
)

// Config defines the structure for exclusion patterns.
type Config struct {
	Exclude []string `json:"exclude"`
}

// LoadConfig loads configuration from a specified path.
// If the file does not exist, it returns a default configuration.
func LoadConfig(configPath string) (*Config, error) {
	// Default configuration with common exclusions
	defaultConfig := &Config{
		Exclude: []string{
			".git/",
			".DS_Store",
			"*.log",
			"go.mod", // Exclude Go module files by default for plugin zipping
			"go.sum",
			"main.go", // Exclude main executable source by default
			"packer/", // Exclude source code directories by default
			"archiver/",
			"config/",
			"cursezip.example.json", // Exclude example config
			".*", // Exclude hidden files by default
		},
	}

	if configPath == "" {
		// Try to load from default path if not specified
		configPath = "cursezip.json"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If config file doesn't exist, return default config
			return defaultConfig, nil
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Merge default exclusions with user-defined ones
	// This simple merge just appends. A more sophisticated merge might handle duplicates or overrides.
	cfg.Exclude = append(defaultConfig.Exclude, cfg.Exclude...)

	return &cfg, nil
}
