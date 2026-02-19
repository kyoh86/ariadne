package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func newUpgradeMCCommand() *cobra.Command {
	var profileFlag string
	var syncNow bool

	cmd := &cobra.Command{
		Use:   "upgrade-mc <mc-version>",
		Short: "Upgrade profile Minecraft version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mcVersion := args[0]
			profileName, err := resolveProfileName(profileFlag)
			if err != nil {
				return err
			}

			profile, err := loadProfile(profileName)
			if err != nil {
				return err
			}

			loader, err := fetchLatestFabricLoader(mcVersion)
			if err != nil {
				return fmt.Errorf("resolve latest fabric loader: %w", err)
			}
			profile.MCVersion = mcVersion
			profile.FabricLoader = loader
			profile.LastUpdatedRFC3339 = time.Now().Format(time.RFC3339)
			if err := saveProfile(profile); err != nil {
				return err
			}

			fmt.Printf("upgraded profile %q to MC %s (Fabric loader %s)\n", profileName, mcVersion, loader)

			if syncNow {
				return syncProfile(profileName)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&profileFlag, "profile", "", "profile name (default: active profile)")
	cmd.Flags().BoolVar(&syncNow, "sync", false, "run sync after upgrade")

	return cmd
}
