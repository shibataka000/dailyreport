package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func runE(ctx context.Context, dir string, query string) error {
	return nil
}

func main() {
	var dir string

	command := &cobra.Command{
		Use:   "dailyreport",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd.Context(), dir, args[0])
		},
		SilenceUsage: true,
	}

	command.Flags().StringVar(&dir, "dir", filepath.Join(os.Getenv("HOME"), "note"), "")

	if err := command.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
