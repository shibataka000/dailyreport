// Package main provides a CLI tool for querying daily reports stored in a specified directory.
// It uses the Cobra library to parse command-line flags for the report directory, date range, and query strings.
// The tool executes queries over the reports within the given date range and prints the results to standard output.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func runE(ctx context.Context, dir string, since time.Time, until time.Time, queries []string) error {
	app := newApplication(newDailyReportRepository(dir))
	for _, query := range queries {
		result, err := app.query(ctx, since, until, query)
		if err != nil {
			return err
		}
		fmt.Println(string(result))
	}
	return nil
}

func main() {
	var (
		dir     string
		since   time.Time
		until   time.Time
		queries []string
	)

	command := &cobra.Command{
		Use:   "dailyreport",
		Short: "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runE(cmd.Context(), dir, since, until, queries)
		},
		SilenceUsage: true,
	}

	command.Flags().StringVar(&dir, "dir", ".", "")
	command.Flags().TimeVar(&since, "since", time.Now().AddDate(0, -1, 0), []string{time.DateOnly}, "")
	command.Flags().TimeVar(&until, "until", time.Now().AddDate(0, 0, 1), []string{time.DateOnly}, "")
	command.Flags().StringSliceVar(&queries, "query", []string{"."}, "")

	if err := command.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
