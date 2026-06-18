package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Camilo-845/typingame/internal/paths"
)

func configPath() (string, error) {
	dir, err := paths.AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

func Load() *Config {
	path, err := configPath()
	if err != nil {
		return DefaultConfig()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig()
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig()
	}

	return cfg.Validate()
}

func Save(cfg *Config) error {
	cfg.Validate()

	dir, err := paths.AppDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
