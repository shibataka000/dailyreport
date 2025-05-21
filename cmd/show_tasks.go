package cmd

import (
	"cmp"
	"fmt"
	"os"
	"time"

	"github.com/shibataka000/dailyreport/dailyreport"
	"github.com/spf13/cobra"
)

// NewShowCommand return new `dailyreport show tasks` sub-command instance.
func NewShowTasksCommand() *cobra.Command {
	var (
		dir      string
		startStr string
		endStr   string
		project  string
	)

	command := &cobra.Command{
		Use:   "tasks",
		Short: "Show tasks in daily report.",
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

			// Union tasks.
			tasks, err = tasks.Union()
			if err != nil {
				return err
			}

			// Sort tasks.
			tasks = tasks.Sort(func(a dailyreport.Task, b dailyreport.Task) int {
				if !a.Project.Equals(b.Project) {
					return cmp.Compare(a.Project.Name, b.Project.Name)
				}
				return cmp.Compare(a.Name, b.Name)
			})

			// Show tasks.
			for _, task := range tasks {
				var completionMark string
				if task.IsCompleted {
					completionMark = "x"
				} else {
					completionMark = " "
				}
				fmt.Printf("- [%s] %.2fh / %.2fh [%s] %s\n", completionMark, task.Estimate.Hours(), task.Actual.Hours(), task.Project.Name, task.Name)
			}

			return nil
		},
	}

	command.Flags().StringVar(&dir, "dir", os.Getenv("DR_DIR"), "Directory where daily report file exists. [$DR_DIR]")
	command.Flags().StringVarP(&startStr, "start-at", "s", os.Getenv("DR_START_AT"), "Start of daily report date range. [$DR_START_AT]")
	command.Flags().StringVarP(&endStr, "end-at", "e", os.Getenv("DR_END_AT"), "End of daily report date range. [$DR_END_AT]")
	command.Flags().StringVar(&project, "project", os.Getenv("DR_PROJECT"), "Show only tasks which project name is this. [$DR_PROJECT]")

	return command
}
