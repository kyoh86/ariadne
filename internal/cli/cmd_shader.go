package cli

import "github.com/spf13/cobra"

func newShaderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shader",
		Short: "Manage shaderpacks",
	}

	cmd.AddCommand(newShaderAddCommand())
	cmd.AddCommand(newShaderListCommand())
	cmd.AddCommand(newShaderDropCommand())
	return cmd
}
