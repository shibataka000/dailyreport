package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		path   string
		report DailyReport
	}{
		{
			path: "./testdata/20250101.md",
			report: DailyReport{
				attendance: Attendance{
					start:     time.Date(2025, 1, 1, 9, 30, 0, 0, time.Local),
					end:       time.Date(2025, 1, 1, 17, 30, 0, 0, time.Local),
					breakTime: 1 * time.Hour,
				},
				tasks: []Task{
					{project: "プロジェクト A", description: "タスク C", estimate: 2*time.Hour + 0*time.Minute, actual: 2*time.Hour + 30*time.Minute, completion: false},
					{project: "プロジェクト A", description: "タスク D", estimate: 2*time.Hour + 0*time.Minute, actual: 1*time.Hour + 30*time.Minute, completion: false},
					{project: "プロジェクト B1", description: "タスク E", estimate: 1*time.Hour + 30*time.Minute, actual: 1*time.Hour + 15*time.Minute, completion: false},
					{project: "プロジェクト B1", description: "タスク F", estimate: 1*time.Hour + 30*time.Minute, actual: 2*time.Hour + 45*time.Minute, completion: false},
				},
			},
		},
		{
			path: "./testdata/20250103.md",
			report: DailyReport{
				attendance: Attendance{
					start:     time.Date(2025, 1, 3, 9, 15, 0, 0, time.Local),
					end:       time.Date(2025, 1, 3, 17, 45, 0, 0, time.Local),
					breakTime: 1*time.Hour + 30*time.Minute,
				},
				tasks: []Task{
					{project: "プロジェクト A", description: "タスク C", estimate: 2*time.Hour + 0*time.Minute, actual: 2*time.Hour + 30*time.Minute, completion: false},
					{project: "プロジェクト A", description: "タスク D", estimate: 2*time.Hour + 0*time.Minute, actual: 1*time.Hour + 30*time.Minute, completion: true},
					{project: "プロジェクト B2", description: "タスク E", estimate: 0*time.Hour + 0*time.Minute, actual: 1*time.Hour + 15*time.Minute, completion: false},
					{project: "プロジェクト B2", description: "タスク F", estimate: 0*time.Hour + 0*time.Minute, actual: 2*time.Hour + 45*time.Minute, completion: false},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			require := require.New(t)
			repository, err := newDailyReportRepository()
			require.NoError(err)
			report, err := repository.read(tt.path)
			require.NoError(err)
			require.Equal(tt.report, report)
		})
	}
}
