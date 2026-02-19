package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const userAgent = "ariadne/0.1 (https://github.com/kyoh86/minecraft-client)"

func resolveVersion(mod ModSpec, mcVersion string) (*modrinthVersion, error) {
	kind := normalizeModKind(mod.Kind)
	if kind == "" {
		kind = "mod"
	}
	if mod.VersionID != "" {
		var ver modrinthVersion
		if err := getJSON("https://api.modrinth.com/v2/version/"+url.PathEscape(mod.VersionID), &ver); err != nil {
			return nil, err
		}
		if kind == "mod" && !contains(ver.Loaders, "fabric") {
			return nil, fmt.Errorf("pinned version %s is not fabric compatible", mod.VersionID)
		}
		if !contains(ver.GameVersions, mcVersion) {
			return nil, fmt.Errorf("pinned version %s does not support Minecraft %s", mod.VersionID, mcVersion)
		}
		return &ver, nil
	}

	gameVersionsJSON := fmt.Sprintf(`["%s"]`, mcVersion)
	endpoint := "https://api.modrinth.com/v2/project/" + url.PathEscape(mod.ProjectID) + "/version" +
		"?game_versions=" + url.QueryEscape(gameVersionsJSON)
	if kind == "mod" {
		loadersJSON := `["fabric"]`
		endpoint += "&loaders=" + url.QueryEscape(loadersJSON)
	}

	var versions []modrinthVersion
	if err := getJSON(endpoint, &versions); err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("no compatible version found for Minecraft %s", mcVersion)
	}
	return &versions[0], nil
}

func pickPrimaryFile(files []modrinthVersionFile) (*modrinthVersionFile, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("version has no files")
	}
	for i := range files {
		if files[i].Primary {
			return &files[i], nil
		}
	}
	return &files[0], nil
}

func fetchProject(ref string) (*modrinthProject, error) {
	var project modrinthProject
	endpoint := "https://api.modrinth.com/v2/project/" + url.PathEscape(strings.TrimSpace(ref))
	if err := getJSON(endpoint, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func fetchLatestFabricLoader(mcVersion string) (string, error) {
	var loaders []fabricLoaderVersion
	endpoint := "https://meta.fabricmc.net/v2/versions/loader/" + url.PathEscape(mcVersion)
	if err := getJSON(endpoint, &loaders); err != nil {
		return "", err
	}
	if len(loaders) == 0 {
		return "", fmt.Errorf("no fabric loader found for MC %s", mcVersion)
	}
	return loaders[0].Loader.Version, nil
}

func getJSON(endpoint string, out any) error {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 8192))
		return fmt.Errorf("request failed: %s: %s", res.Status, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(res.Body).Decode(out)
}
