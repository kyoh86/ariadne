package cli

import "github.com/spf13/cobra"

func newModCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "Manage mods",
	}

	cmd.AddCommand(newModAddCommand())
	cmd.AddCommand(newModDropCommand())
	cmd.AddCommand(newModListCommand())
	cmd.AddCommand(newModSearchCommand())
	return cmd
}
