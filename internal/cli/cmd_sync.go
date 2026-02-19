package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func newSyncCommand() *cobra.Command {
	var profileFlag string

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Resolve and sync mods to game directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, err := resolveProfileName(profileFlag)
			if err != nil {
				return err
			}
			return syncProfile(profileName)
		},
	}

	cmd.Flags().StringVar(&profileFlag, "profile", "", "profile name (default: active profile)")
	return cmd
}

func syncProfile(profileName string) error {
	if profileName == "" {
		return fmt.Errorf("profile is required")
	}

	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}
	if profile.FabricLoader == "" || profile.FabricLoader == "latest" {
		loader, err := fetchLatestFabricLoader(profile.MCVersion)
		if err != nil {
			return fmt.Errorf("resolve latest fabric loader: %w", err)
		}
		profile.FabricLoader = loader
		profile.LastUpdatedRFC3339 = time.Now().Format(time.RFC3339)
		if err := saveProfile(profile); err != nil {
			return err
		}
	}

	modsDir := filepath.Join(profile.GameDir, "mods")
	shaderpacksDir := filepath.Join(profile.GameDir, "shaderpacks")
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		return fmt.Errorf("create mods dir: %w", err)
	}
	if err := os.MkdirAll(shaderpacksDir, 0o755); err != nil {
		return fmt.Errorf("create shaderpacks dir: %w", err)
	}

	oldLock, _ := loadLock(profile.Name)
	newLock := Lock{
		ProfileName:        profile.Name,
		MCVersion:          profile.MCVersion,
		FabricLoader:       profile.FabricLoader,
		GeneratedAtRFC3339: time.Now().Format(time.RFC3339),
		Entries:            map[string]LockEntry{},
	}

	for _, mod := range profile.Mods {
		kind := normalizeModKind(mod.Kind)
		if kind == "" {
			kind = "mod"
		}
		version, err := resolveVersion(mod, profile.MCVersion)
		if err != nil {
			return fmt.Errorf("resolve version for %s: %w", mod.Slug, err)
		}
		file, err := pickPrimaryFile(version.Files)
		if err != nil {
			return fmt.Errorf("pick file for %s: %w", mod.Slug, err)
		}
		targetDir := targetDirForKind(kind)
		targetBase := modsDir
		if targetDir == "shaderpacks" {
			targetBase = shaderpacksDir
		}
		dst := filepath.Join(targetBase, file.Filename)
		if err := downloadWithOptionalSHA512(file.URL, dst, file.Hashes["sha512"]); err != nil {
			return fmt.Errorf("download %s: %w", mod.Slug, err)
		}

		lockKey := kind + ":" + mod.ProjectID
		newLock.Entries[lockKey] = LockEntry{
			ProjectID:  mod.ProjectID,
			Slug:       mod.Slug,
			Kind:       kind,
			VersionID:  version.ID,
			TargetDir:  targetDir,
			FileName:   file.Filename,
			FileSHA512: file.Hashes["sha512"],
		}
		fmt.Printf("synced %-24s (%s) -> %s/%s\n", mod.Slug, kind, targetDir, file.Filename)
	}

	if oldLock != nil {
		for lockKey, old := range oldLock.Entries {
			newEntry, ok := newLock.Entries[lockKey]
			if ok && newEntry.FileName == old.FileName {
				continue
			}
			staleDir := old.TargetDir
			if staleDir == "" {
				staleDir = targetDirForKind(old.Kind)
			}
			staleBase := modsDir
			if staleDir == "shaderpacks" {
				staleBase = shaderpacksDir
			}
			stale := filepath.Join(staleBase, old.FileName)
			if err := os.Remove(stale); err == nil {
				fmt.Printf("removed stale file: %s/%s\n", staleDir, old.FileName)
			}
		}
	}

	if err := saveLock(newLock); err != nil {
		return err
	}

	fmt.Printf("sync completed: %d items (%s / fabric-loader %s)\n", len(newLock.Entries), profile.MCVersion, profile.FabricLoader)
	return nil
}
