package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newShaderDropCommand() *cobra.Command {
	var profileFlag string

	cmd := &cobra.Command{
		Use:   "drop <project>",
		Short: "Drop a shaderpack from profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, err := resolveProfileName(profileFlag)
			if err != nil {
				return err
			}
			projectRef := strings.TrimSpace(args[0])
			return dropItemFromProfile(profileName, projectRef, "shader")
		},
	}

	cmd.Flags().StringVar(&profileFlag, "profile", "", "profile name (default: active profile)")
	return cmd
}
