package cli

import "github.com/spf13/cobra"

func newModListCommand() *cobra.Command {
	var profileFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List mods in profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, err := resolveProfileName(profileFlag)
			if err != nil {
				return err
			}
			return listItemsInProfile(profileName, "mod")
		},
	}

	cmd.Flags().StringVar(&profileFlag, "profile", "", "profile name (default: active profile)")
	return cmd
}
