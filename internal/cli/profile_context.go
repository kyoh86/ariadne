package cli

import (
	"fmt"
	"strings"
)

func resolveProfileName(maybe string) (string, error) {
	name := strings.TrimSpace(maybe)
	if name != "" {
		return name, nil
	}
	active, err := configuredActiveProfile()
	if err != nil {
		return "", err
	}
	active = strings.TrimSpace(active)
	if active == "" {
		return "", fmt.Errorf("profile is required; set one with `ariadne profile use <name>`")
	}
	return active, nil
}
