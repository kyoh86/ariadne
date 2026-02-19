package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newProfileDropCommand() *cobra.Command {
	var deleteGameDir bool

	cmd := &cobra.Command{
		Use:   "drop <name>",
		Short: "Drop profile definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			profile, err := loadProfile(name)
			if err != nil {
				return err
			}

			if err := deleteProfileArtifacts(name); err != nil {
				return err
			}
			if cfg, err := loadConfig(); err == nil && cfg.ActiveProfile == name {
				cfg.ActiveProfile = ""
				_ = saveConfig(cfg)
			}

			launcherDir, err := configuredLauncherDir()
			if err == nil {
				absLauncherDir, err := filepath.Abs(launcherDir)
				if err == nil {
					_ = removeLauncherProfile(absLauncherDir, name)
				}
			}

			if deleteGameDir {
				if err := os.RemoveAll(profile.GameDir); err != nil {
					return fmt.Errorf("delete game-dir: %w", err)
				}
				fmt.Printf("dropped profile %q and deleted game-dir: %s\n", name, profile.GameDir)
				return nil
			}

			fmt.Printf("dropped profile %q\n", name)
			fmt.Printf("game-dir remains: %s\n", profile.GameDir)
			return nil
		},
	}

	cmd.Flags().BoolVar(&deleteGameDir, "delete-game-dir", false, "also delete game directory")
	return cmd
}
