package cli

import "github.com/spf13/cobra"

// listCmd is an alias for subjects.
func (a *App) listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all CS subject guides (alias for subjects)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			subjects, err := a.client.Subjects(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			if a.limit > 0 && len(subjects) > a.limit {
				subjects = subjects[:a.limit]
			}
			return a.renderOrEmpty(subjects, len(subjects))
		},
	}
}

// topicCmd fetches one subject by slug.
func (a *App) topicCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "topic <slug>",
		Short: "Show details for a CS topic by slug (e.g. programming)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			subjects, err := a.client.Subjects(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			for _, s := range subjects {
				if s.Slug == args[0] {
					return a.render([]*subjectSlice{
						{Rank: s.Rank, Slug: s.Slug, Title: s.Title, URL: s.URL, BookURL: s.BookURL},
					})
				}
			}
			return codeError(exitNoData, nil)
		},
	}
}

// subjectSlice is render-friendly (same fields as Subject).
type subjectSlice struct {
	Rank    int    `json:"rank"     csv:"rank"     tsv:"rank"`
	Slug    string `json:"slug"     csv:"slug"     tsv:"slug"`
	Title   string `json:"title"    csv:"title"    tsv:"title"`
	URL     string `json:"url"      csv:"url"      tsv:"url"`
	BookURL string `json:"book_url" csv:"book_url" tsv:"book_url"`
}

// infoCmd prints site stats.
func (a *App) infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print site stats (subject count)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := a.client.SiteInfo(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			return a.render(info)
		},
	}
}
