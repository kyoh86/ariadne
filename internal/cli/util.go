package cli

import "strings"

func contains(values []string, needle string) bool {
	for _, v := range values {
		if v == needle {
			return true
		}
	}
	return false
}

func normalizeModKind(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "", "mod":
		return "mod"
	case "shader", "shaderpack":
		return "shader"
	default:
		return ""
	}
}

func targetDirForKind(kind string) string {
	switch normalizeModKind(kind) {
	case "shader":
		return "shaderpacks"
	default:
		return "mods"
	}
}
