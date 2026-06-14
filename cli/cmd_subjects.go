package cli

import "github.com/spf13/cobra"

func (a *App) subjectsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "subjects",
		Short: "List all CS subject guides from Teach Yourself CS",
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
