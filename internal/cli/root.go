package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ariadne",
		Short: "Manage Fabric + Modrinth based Minecraft mod profiles",
	}

	cmd.AddCommand(newProfileCommand())
	cmd.AddCommand(newModCommand())
	cmd.AddCommand(newShaderCommand())
	cmd.AddCommand(newSyncCommand())
	cmd.AddCommand(newUpgradeMCCommand())
	cmd.AddCommand(newConfigCommand())

	return cmd
}
