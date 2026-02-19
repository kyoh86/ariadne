package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const appDirName = "ariadne"

func ensureDataDirs() error {
	if err := os.MkdirAll(profileDir(), 0o755); err != nil {
		return fmt.Errorf("create profile dir: %w", err)
	}
	if err := os.MkdirAll(lockDir(), 0o755); err != nil {
		return fmt.Errorf("create lock dir: %w", err)
	}
	if err := os.MkdirAll(configDir(), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return nil
}

func loadProfile(name string) (Profile, error) {
	var profile Profile
	path := profilePath(name)
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return profile, fmt.Errorf("profile %q not found: run `profile create` first", name)
		}
		return profile, err
	}
	if err := json.Unmarshal(b, &profile); err != nil {
		return profile, fmt.Errorf("parse profile: %w", err)
	}
	return profile, nil
}

func saveProfile(profile Profile) error {
	if err := ensureDataDirs(); err != nil {
		return err
	}
	return writeJSON(profilePath(profile.Name), profile)
}

func listProfiles() ([]Profile, error) {
	entries, err := os.ReadDir(profileDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []Profile{}, nil
		}
		return nil, err
	}
	profiles := make([]Profile, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		p, err := loadProfile(name)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})
	return profiles, nil
}

func deleteProfileArtifacts(name string) error {
	if err := os.Remove(profilePath(name)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(lockPath(name)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func loadLock(profileName string) (*Lock, error) {
	b, err := os.ReadFile(lockPath(profileName))
	if err != nil {
		return nil, err
	}
	var lock Lock
	if err := json.Unmarshal(b, &lock); err != nil {
		return nil, err
	}
	return &lock, nil
}

func saveLock(lock Lock) error {
	if err := ensureDataDirs(); err != nil {
		return err
	}
	return writeJSON(lockPath(lock.ProfileName), lock)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

func dataDir() string {
	if base := strings.TrimSpace(os.Getenv("XDG_DATA_HOME")); base != "" {
		return filepath.Join(base, appDirName)
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".", ".local", "share", appDirName)
	}
	return filepath.Join(home, ".local", "share", appDirName)
}

func configDir() string {
	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		return filepath.Join(".", ".config", appDirName)
	}
	return filepath.Join(base, appDirName)
}

func profileDir() string {
	return filepath.Join(dataDir(), "profiles")
}

func lockDir() string {
	return filepath.Join(dataDir(), "locks")
}

func profilePath(name string) string {
	return filepath.Join(profileDir(), name+".json")
}

func lockPath(name string) string {
	return filepath.Join(lockDir(), name+".json")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}
