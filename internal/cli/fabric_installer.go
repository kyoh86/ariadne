package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type fabricInstallerVersion struct {
	Version string `json:"version"`
}

func ensureFabricLoaderInitialized(launcherDir, mcVersion, loaderVersion string) error {
	if err := os.MkdirAll(launcherDir, 0o755); err != nil {
		return fmt.Errorf("create launcher dir: %w", err)
	}
	cacheDir := filepath.Join(dataDir(), "cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	installerVersion, err := fetchLatestFabricInstallerVersion()
	if err != nil {
		return err
	}
	installerJarPath := filepath.Join(cacheDir, fmt.Sprintf("fabric-installer-%s.jar", installerVersion))
	if _, err := os.Stat(installerJarPath); os.IsNotExist(err) {
		url := fmt.Sprintf("https://maven.fabricmc.net/net/fabricmc/fabric-installer/%s/fabric-installer-%s.jar", installerVersion, installerVersion)
		if err := downloadWithOptionalSHA512(url, installerJarPath, ""); err != nil {
			return fmt.Errorf("download fabric installer: %w", err)
		}
	} else if err != nil {
		return err
	}

	javaPath, err := exec.LookPath("java")
	if err != nil {
		return fmt.Errorf("java not found in PATH")
	}

	args := []string{
		"-jar", installerJarPath,
		"client",
		"-dir", launcherDir,
		"-mcversion", mcVersion,
		"-loader", loaderVersion,
		"-noprofile",
	}
	out, err := exec.Command(javaPath, args...).CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("fabric installer failed: %s", msg)
	}
	versionID := fabricVersionID(mcVersion, loaderVersion)
	versionJSON := filepath.Join(launcherDir, "versions", versionID, versionID+".json")
	if _, err := os.Stat(versionJSON); err != nil {
		return fmt.Errorf("fabric installer completed but version metadata not found: %s", versionJSON)
	}
	return nil
}

func fetchLatestFabricInstallerVersion() (string, error) {
	var versions []fabricInstallerVersion
	if err := getJSON("https://meta.fabricmc.net/v2/versions/installer", &versions); err != nil {
		return "", fmt.Errorf("resolve latest fabric installer: %w", err)
	}
	if len(versions) == 0 || versions[0].Version == "" {
		return "", fmt.Errorf("fabric installer version not found")
	}
	return versions[0].Version, nil
}

func fabricVersionID(mcVersion, loaderVersion string) string {
	return fmt.Sprintf("fabric-loader-%s-%s", loaderVersion, mcVersion)
}
