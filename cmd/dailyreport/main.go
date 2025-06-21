// Package main is the entrypoint for the dailyreport CLI tool.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shibataka000/dailyreport"
	"github.com/spf13/cobra"
)

func main() {
	var from, to, dir string

	rootCmd := &cobra.Command{
		Use:   "dailyreport",
		Short: "日報集計CLIツール",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			repo, err := dailyreport.NewReportRepository(ctx, dir)
			if err != nil {
				return err
			}
			app, err := dailyreport.NewApplicationService(ctx, repo)
			if err != nil {
				return err
			}
			fromTime, err := time.Parse("20060102", from)
			if err != nil {
				return err
			}
			toTime, err := time.Parse("20060102", to)
			if err != nil {
				return err
			}
			result, err := app.Aggregate(ctx, fromTime, toTime)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}

	rootCmd.Flags().StringVar(&from, "from", time.Now().Format("20060102"), "集計期間の開始日 (YYYYMMDD)")
	rootCmd.Flags().StringVar(&to, "to", time.Now().Format("20060102"), "集計期間の終了日 (YYYYMMDD)")
	rootCmd.Flags().StringVar(&dir, "dir", "", "日報が格納されているディレクトリ")
	if err := rootCmd.MarkFlagRequired("dir"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
