package cli

type Profile struct {
	Name               string    `json:"name"`
	MCVersion          string    `json:"mcVersion"`
	FabricLoader       string    `json:"fabricLoader"`
	GameDir            string    `json:"gameDir"`
	Mods               []ModSpec `json:"mods"`
	LastUpdatedRFC3339 string    `json:"lastUpdated"`
}

type ModSpec struct {
	ProjectID string `json:"projectId"`
	Slug      string `json:"slug"`
	Kind      string `json:"kind,omitempty"` // mod | shader
	VersionID string `json:"versionId,omitempty"`
}

type Lock struct {
	ProfileName        string               `json:"profileName"`
	MCVersion          string               `json:"mcVersion"`
	FabricLoader       string               `json:"fabricLoader"`
	GeneratedAtRFC3339 string               `json:"generatedAt"`
	Entries            map[string]LockEntry `json:"entries"`
}

type LockEntry struct {
	ProjectID  string `json:"projectId"`
	Slug       string `json:"slug"`
	Kind       string `json:"kind,omitempty"` // mod | shader
	VersionID  string `json:"versionId"`
	TargetDir  string `json:"targetDir,omitempty"` // mods | shaderpacks
	FileName   string `json:"fileName"`
	FileSHA512 string `json:"fileSha512,omitempty"`
}

type modrinthProject struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

type modrinthVersion struct {
	ID           string                `json:"id"`
	GameVersions []string              `json:"game_versions"`
	Loaders      []string              `json:"loaders"`
	Files        []modrinthVersionFile `json:"files"`
}

type modrinthVersionFile struct {
	Filename string            `json:"filename"`
	URL      string            `json:"url"`
	Primary  bool              `json:"primary"`
	Hashes   map[string]string `json:"hashes"`
}

type fabricLoaderVersion struct {
	Loader struct {
		Version string `json:"version"`
	} `json:"loader"`
}
