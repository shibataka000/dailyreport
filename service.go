package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
	"strings"
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

func jq(ctx context.Context, output JQOutput, filter string) (string, error) {
	temp, err := os.CreateTemp("", "*")
	if err != nil {
		return "", err
	}
	defer os.Remove(temp.Name()) // nolint:errcheck

	b1, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	if _, err := temp.Write(b1); err != nil {
		return "", err
	}
	if err = temp.Close(); err != nil {
		return "", err
	}

	b2, err := exec.CommandContext(ctx, "jq", "-r", filter, temp.Name()).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b2)), nil
}
