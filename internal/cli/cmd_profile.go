package cli

import "github.com/spf13/cobra"

func newProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}

	cmd.AddCommand(newInitProfileCommand())
	cmd.AddCommand(newProfileListCommand())
	cmd.AddCommand(newProfileUseCommand())
	cmd.AddCommand(newProfileDropCommand())
	return cmd
}
