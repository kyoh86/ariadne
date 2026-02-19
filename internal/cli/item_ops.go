package cli

import (
	"fmt"
	"strings"
	"time"
)

func addItemToProfile(profileName, projectRef, versionID, kind string) error {
	kind = normalizeModKind(kind)
	if kind == "" {
		return fmt.Errorf("invalid kind: %s", kind)
	}

	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}

	project, err := fetchProject(projectRef)
	if err != nil {
		return fmt.Errorf("resolve project %q: %w", projectRef, err)
	}

	replaced := false
	for i := range profile.Mods {
		if profile.Mods[i].ProjectID == project.ID && normalizeModKind(profile.Mods[i].Kind) == kind {
			profile.Mods[i].Slug = project.Slug
			profile.Mods[i].Kind = kind
			profile.Mods[i].VersionID = strings.TrimSpace(versionID)
			replaced = true
			break
		}
	}
	if !replaced {
		profile.Mods = append(profile.Mods, ModSpec{
			ProjectID: project.ID,
			Slug:      project.Slug,
			Kind:      kind,
			VersionID: strings.TrimSpace(versionID),
		})
	}
	profile.LastUpdatedRFC3339 = time.Now().Format(time.RFC3339)
	if err := saveProfile(profile); err != nil {
		return err
	}

	if versionID == "" {
		fmt.Printf("added %s %q to profile %q\n", kind, project.Slug, profileName)
	} else {
		fmt.Printf("added %s %q to profile %q (pinned version %s)\n", kind, project.Slug, profileName, versionID)
	}
	return nil
}

func listItemsInProfile(profileName, kind string) error {
	kind = normalizeModKind(kind)
	if kind == "" {
		return fmt.Errorf("invalid kind: %s", kind)
	}
	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}
	count := 0
	for _, m := range profile.Mods {
		mKind := normalizeModKind(m.Kind)
		if mKind != kind {
			continue
		}
		pinned := m.VersionID
		if pinned == "" {
			pinned = "-"
		}
		fmt.Printf("%s\t%s\tversion-id=%s\n", m.Slug, m.ProjectID, pinned)
		count++
	}
	if count == 0 {
		fmt.Printf("no %ss in profile %q\n", kind, profile.Name)
	}
	return nil
}

func dropItemFromProfile(profileName, projectRef, kind string) error {
	kind = normalizeModKind(kind)
	if kind == "" {
		return fmt.Errorf("invalid kind: %s", kind)
	}
	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}

	ref := strings.TrimSpace(projectRef)
	newMods := make([]ModSpec, 0, len(profile.Mods))
	removed := false
	for _, m := range profile.Mods {
		if normalizeModKind(m.Kind) == kind && (m.Slug == ref || m.ProjectID == ref) {
			removed = true
			continue
		}
		newMods = append(newMods, m)
	}
	if !removed {
		return fmt.Errorf("%s %q not found in profile %q", kind, ref, profileName)
	}

	profile.Mods = newMods
	profile.LastUpdatedRFC3339 = time.Now().Format(time.RFC3339)
	if err := saveProfile(profile); err != nil {
		return err
	}
	fmt.Printf("dropped %s %q from profile %q\n", kind, ref, profileName)
	return nil
}
