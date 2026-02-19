package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

type modrinthSearchResponse struct {
	Hits []modrinthSearchHit `json:"hits"`
}

type modrinthSearchHit struct {
	ProjectID   string `json:"project_id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func newModSearchCommand() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search mods on Modrinth",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.TrimSpace(args[0])
			if query == "" {
				return fmt.Errorf("query is required")
			}
			if limit <= 0 {
				return fmt.Errorf("limit must be > 0")
			}

			endpoint := "https://api.modrinth.com/v2/search?query=" + url.QueryEscape(query) +
				"&limit=" + url.QueryEscape(fmt.Sprintf("%d", limit)) +
				"&index=" + url.QueryEscape("relevance")
			var res modrinthSearchResponse
			if err := getJSON(endpoint, &res); err != nil {
				return err
			}

			if len(res.Hits) == 0 {
				fmt.Println("no results")
				return nil
			}
			for _, hit := range res.Hits {
				desc := strings.ReplaceAll(hit.Description, "\n", " ")
				fmt.Printf("%s\t%s\t%s\n", hit.Slug, hit.ProjectID, desc)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "maximum number of results")
	return cmd
}
