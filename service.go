package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
)

func aggregate(reports []DailyReport) AggregatedReport {
	all := []Task{}
	for _, report := range reports {
		all = append(all, report.Tasks...)
	}

	aggregated := []Task{}
	for _, task := range all {
		i := slices.IndexFunc(aggregated, func(t Task) bool {
			return task.Project == t.Project && task.Description == t.Description
		})
		if i == -1 {
			aggregated = append(aggregated, task)
		} else {
			aggregated[i].Estimated += task.Estimated
			aggregated[i].Actual += task.Actual
			aggregated[i].Completed = aggregated[i].Completed || task.Completed
		}
	}

	return AggregatedReport{
		Tasks: aggregated,
	}
}

func jq(ctx context.Context, output JQOutput, filter string) ([]byte, error) {
	temp, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(temp.Name())

	b, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	if _, err := temp.Write(b); err != nil {
		return nil, err
	}
	if err = temp.Close(); err != nil {
		return nil, err
	}

	return exec.CommandContext(ctx, "jq", filter, temp.Name()).Output()
}
