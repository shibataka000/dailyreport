package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DailyReportRepository struct{}

func newDailyReportRepository() (*DailyReportRepository, error) {
	return &DailyReportRepository{}, nil
}

func (r *DailyReportRepository) read(path string) (DailyReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DailyReport{}, err
	}

	t, err := time.Parse("20060102.md", filepath.Base(path))
	if err != nil {
		return DailyReport{}, err
	}

	return unmarshal(t, data)
}

func unmarshal(t time.Time, data []byte) (DailyReport, error) {
	var (
		projectPattern = regexp.MustCompile(`^- \[.\] (.+)$`)
		taskPattern    = regexp.MustCompile(`^  - \[(.?)\] ([^\s]+)/([^\s]+) (.+)$`)

		report         DailyReport
		currentProject string
	)

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "- 始業 "):
			d, err := parseDuration("- 始業 15:04", line)
			if err != nil {
				return report, err
			}
			report.attendance.start = t.Add(d)
		case strings.HasPrefix(line, "- 終業 "):
			d, err := parseDuration("- 終業 15:04", line)
			if err != nil {
				return report, err
			}
			report.attendance.end = t.Add(d)
		case strings.HasPrefix(line, "- 休憩 "):
			d, err := parseDuration("- 休憩 15:04", line)
			if err != nil {
				return report, err
			}
			report.attendance.breakTime = d
		case projectPattern.MatchString(line):
			currentProject = projectPattern.FindStringSubmatch(line)[1]
		case taskPattern.MatchString(line):
			matches := taskPattern.FindStringSubmatch(line)
			estimate, err := time.ParseDuration(matches[2])
			if err != nil {
				return report, err
			}
			actual, err := time.ParseDuration(matches[3])
			if err != nil {
				return report, err
			}
			report.tasks = append(report.tasks, Task{
				project:     currentProject,
				description: matches[4],
				estimate:    estimate,
				actual:      actual,
				completion:  matches[1] == "x",
			})
		}
		if line == "---" {
			break
		}
	}

	return report, nil
}

func parseDuration(layout string, value string) (time.Duration, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return 0, err
	}
	return t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())), nil
}
