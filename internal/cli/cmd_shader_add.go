package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newShaderAddCommand() *cobra.Command {
	var profileFlag, versionID string

	cmd := &cobra.Command{
		Use:   "add <project>",
		Short: "Add a shaderpack project to profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRef := strings.TrimSpace(args[0])
			profileName, err := resolveProfileName(profileFlag)
			if err != nil {
				return err
			}
			return addItemToProfile(profileName, projectRef, versionID, "shader")
		},
	}

	cmd.Flags().StringVar(&profileFlag, "profile", "", "profile name (default: active profile)")
	cmd.Flags().StringVar(&versionID, "version-id", "", "modrinth version id")
	return cmd
}
