package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newInitProfileCommand() *cobra.Command {
	var nameFlag, gameDir, fabricLoader string

	cmd := &cobra.Command{
		Use:   "create <mc-version>",
		Short: "Create a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mcVersion := strings.TrimSpace(args[0])
			name := strings.TrimSpace(nameFlag)
			if name == "" {
				name = "fabricmc-" + mcVersion
			}
			if name == "" || mcVersion == "" {
				return fmt.Errorf("name and mc-version are required")
			}
			effectiveGameDir := strings.TrimSpace(gameDir)
			if effectiveGameDir == "" {
				root, err := configuredGameDirRoot()
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
				effectiveGameDir = filepath.Join(root, name)
			}

			loader := strings.TrimSpace(fabricLoader)
			if loader == "" || loader == "latest" {
				latest, err := fetchLatestFabricLoader(mcVersion)
				if err != nil {
					return fmt.Errorf("resolve latest fabric loader: %w", err)
				}
				loader = latest
			}

			absGameDir, err := filepath.Abs(effectiveGameDir)
			if err != nil {
				return fmt.Errorf("resolve game-dir: %w", err)
			}
			if err := os.MkdirAll(filepath.Join(absGameDir, "mods"), 0o755); err != nil {
				return fmt.Errorf("create mods dir: %w", err)
			}
			launcherDir, err := configuredLauncherDir()
			if err != nil {
				return fmt.Errorf("load launcher-dir config: %w", err)
			}
			absLauncherDir, err := filepath.Abs(launcherDir)
			if err != nil {
				return fmt.Errorf("resolve launcher-dir: %w", err)
			}
			if err := ensureFabricLoaderInitialized(absLauncherDir, mcVersion, loader); err != nil {
				return err
			}
			if err := upsertLauncherProfile(absLauncherDir, name, absGameDir, mcVersion, loader); err != nil {
				return err
			}

			profile := Profile{
				Name:               name,
				MCVersion:          mcVersion,
				FabricLoader:       loader,
				GameDir:            absGameDir,
				Mods:               []ModSpec{},
				LastUpdatedRFC3339: time.Now().Format(time.RFC3339),
			}
			if err := saveProfile(profile); err != nil {
				return err
			}

			fmt.Printf("initialized profile %q (MC %s, Fabric loader %s)\n", name, mcVersion, loader)
			fmt.Printf("game-dir: %s\n", absGameDir)
			versionID := fabricVersionID(mcVersion, loader)
			fmt.Printf("fabric-loader initialized: %s/versions/%s\n", absLauncherDir, versionID)
			fmt.Printf("launcher profile updated: %s/launcher_profiles.json\n", absLauncherDir)
			return nil
		},
	}

	cmd.Flags().StringVar(&nameFlag, "name", "", "profile name (default: fabricmc-<mc-version>)")
	cmd.Flags().StringVar(&gameDir, "game-dir", "", "game directory path (default: <configured-root>/<name>)")
	cmd.Flags().StringVar(&fabricLoader, "fabric-loader-version", "latest", "fabric loader version or latest")

	return cmd
}
