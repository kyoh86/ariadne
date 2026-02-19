package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func upsertLauncherProfile(launcherDir, profileName, gameDir, mcVersion, loaderVersion string) error {
	if err := os.MkdirAll(launcherDir, 0o755); err != nil {
		return fmt.Errorf("create launcher dir: %w", err)
	}
	path := filepath.Join(launcherDir, "launcher_profiles.json")

	root := map[string]any{}
	if b, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(b, &root); err != nil {
			return fmt.Errorf("parse launcher_profiles.json: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	profiles, ok := root["profiles"].(map[string]any)
	if !ok || profiles == nil {
		profiles = map[string]any{}
	}

	now := time.Now().Format(time.RFC3339)
	lastVersionID := fabricVersionID(mcVersion, loaderVersion)

	entry := map[string]any{
		"name":          profileName,
		"type":          "custom",
		"created":       now,
		"lastUsed":      now,
		"icon":          "Furnace",
		"lastVersionId": lastVersionID,
		"gameDir":       gameDir,
	}
	if existing, ok := profiles[profileName].(map[string]any); ok {
		if created, ok := existing["created"]; ok {
			entry["created"] = created
		}
		for k, v := range existing {
			if _, keep := entry[k]; !keep {
				entry[k] = v
			}
		}
	}

	profiles[profileName] = entry
	root["profiles"] = profiles
	if _, ok := root["selectedProfile"]; !ok {
		root["selectedProfile"] = profileName
	}

	return writeJSON(path, root)
}

func removeLauncherProfile(launcherDir, profileName string) error {
	path := filepath.Join(launcherDir, "launcher_profiles.json")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	root := map[string]any{}
	if err := json.Unmarshal(b, &root); err != nil {
		return fmt.Errorf("parse launcher_profiles.json: %w", err)
	}

	profiles, ok := root["profiles"].(map[string]any)
	if !ok || profiles == nil {
		return nil
	}
	delete(profiles, profileName)
	root["profiles"] = profiles
	if selected, ok := root["selectedProfile"].(string); ok && selected == profileName {
		root["selectedProfile"] = ""
	}
	return writeJSON(path, root)
}
