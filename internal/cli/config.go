package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GameDirRoot   string `json:"gameDirRoot"`
	LauncherDir   string `json:"launcherDir"`
	ActiveProfile string `json:"activeProfile,omitempty"`
}

func loadConfig() (Config, error) {
	var cfg Config
	b, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			cfg.GameDirRoot = defaultGameDirRoot()
			cfg.LauncherDir = cfg.GameDirRoot
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	if cfg.GameDirRoot == "" {
		cfg.GameDirRoot = defaultGameDirRoot()
	}
	if cfg.LauncherDir == "" {
		cfg.LauncherDir = cfg.GameDirRoot
	}
	return cfg, nil
}

func saveConfig(cfg Config) error {
	if err := ensureDataDirs(); err != nil {
		return err
	}
	return writeJSON(configPath(), cfg)
}

func configuredGameDirRoot() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	return cfg.GameDirRoot, nil
}

func configuredLauncherDir() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	return cfg.LauncherDir, nil
}

func configuredActiveProfile() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	return cfg.ActiveProfile, nil
}

func defaultGameDirRoot() string {
	if root := os.Getenv("ARIADNE_GAME_DIR_ROOT"); root != "" {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "./minecraft"
	}
	return filepath.Join(home, ".minecraft")
}
