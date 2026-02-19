package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newProfileUseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Set active profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if _, err := loadProfile(name); err != nil {
				return err
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			cfg.ActiveProfile = name
			if err := saveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("active profile: %s\n", name)
			return nil
		},
	}
}
