package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newProfileListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := listProfiles()
			if err != nil {
				return err
			}
			active, _ := configuredActiveProfile()
			if len(profiles) == 0 {
				fmt.Println("no profiles")
				return nil
			}
			for _, p := range profiles {
				marker := " "
				if p.Name == active {
					marker = "*"
				}
				fmt.Printf("%s %s\tmc=%s\tgame-dir=%s\n", marker, p.Name, p.MCVersion, p.GameDir)
			}
			return nil
		},
	}
}
