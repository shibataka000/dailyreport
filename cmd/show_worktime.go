package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/shibataka000/dailyreport/dailyreport"
	"github.com/spf13/cobra"
)

// NewShowCommand return new `dailyreport show worktime` sub-command instance.
func NewShowWorkTimeCommand() *cobra.Command {
	var (
		dir      string
		startStr string
		endStr   string
		project  string
	)

	command := &cobra.Command{
		Use:   "worktime",
		Short: "Show work time in daily report.",
		RunE: func(_ *cobra.Command, _ []string) error {
			// Read daily reports.
			app := dailyreport.NewApplicationService(dailyreport.NewRepository(dir))
			start, err := time.Parse("20060102", startStr)
			if err != nil {
				return err
			}
			end, err := time.Parse("20060102", endStr)
			if err != nil {
				return err
			}
			reports, err := app.Read(start, end)
			if err != nil {
				return err
			}

			// Get tasks.
			tasks, err := reports.Tasks()
			if err != nil {
				return err
			}

			// Filter tasks.
			if project != "" {
				prj := dailyreport.NewProject(project)
				tasks = tasks.Filter(func(t dailyreport.Task) bool {
					return t.Project.Equals(prj)
				})
			}

			// Show worktime.
			fmt.Printf("Work Time\t\t%.2fh\n", reports.WorkTimes().Duration().Hours())
			fmt.Printf("Tasks (Estimated)\t%.2fh\n", tasks.Estimate().Hours())
			fmt.Printf("Tasks (Actual)\t\t%.2fh\n", tasks.Actual().Hours())

			return nil
		},
	}

	command.Flags().StringVar(&dir, "dir", os.Getenv("DR_DIR"), "Directory where daily report file exists. [$DR_DIR]")
	command.Flags().StringVarP(&startStr, "start-at", "s", os.Getenv("DR_START_AT"), "Start of daily report date range. [$DR_START_AT]")
	command.Flags().StringVarP(&endStr, "end-at", "e", os.Getenv("DR_END_AT"), "End of daily report date range. [$DR_END_AT]")
	command.Flags().StringVar(&project, "project", os.Getenv("DR_PROJECT"), "Show only tasks which project name is this. [$DR_PROJECT]")

	return command
}
