package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

type LLMConfig struct {
	Endpoint string   `toml:"api_base_url"`
	Models   []string `toml:"models"`
}

type RevlyConfig struct {
	LLM LLMConfig `toml:"llm"`
}

var (
	config     RevlyConfig
	configErr  error
	loadedOnce sync.Once
)

// GetConfig loads and returns the parsed Revly config. Only loads once (singleton).
func GetConfig() (RevlyConfig, error) {
	loadedOnce.Do(func() {
		pathsToTry := []string{}

		// Current working directory
		cwd, err := os.Getwd()
		if err == nil {
			pathsToTry = append(pathsToTry, filepath.Join(cwd, "revly.config.toml"))
		}

		// ~/.revly/config.toml and ~/revly.config.toml
		homeDir, err := os.UserHomeDir()
		if err != nil {
			configErr = fmt.Errorf("unable to find user home directory: %w", err)
			return
		}
		pathsToTry = append(pathsToTry,
			filepath.Join(homeDir, ".revly", "config.toml"),
			filepath.Join(homeDir, "revly.config.toml"),
		)

		// Try loading from each path
		for _, path := range pathsToTry {
			if _, err := os.Stat(path); err == nil {
				if _, err := toml.DecodeFile(path, &config); err != nil {
					configErr = fmt.Errorf("failed to parse revly config at %s: %w", path, err)
				} else {
					configErr = nil
				}
				return
			}
		}

		configErr = fmt.Errorf("revly config not found in: %v", pathsToTry)
	})

	return config, configErr
}