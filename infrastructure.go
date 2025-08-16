package main

import (
	"bufio"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

type DailyReportRepository struct {
	dir string
}

func newDailyReportRepository(dir string) *DailyReportRepository {
	return &DailyReportRepository{
		dir: dir,
	}
}

func (r *DailyReportRepository) list(since time.Time, until time.Time) ([]DailyReport, error) {
	reports := []DailyReport{}
	err := filepath.Walk(r.dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isDailyReportFile(path, info) {
			report, err := read(path)
			if err != nil {
				return err
			}
			reports = append(reports, report)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sort(filter(reports, since, until)), nil
}

func isDailyReportFile(path string, info fs.FileInfo) bool {
	_, err := time.Parse("20060102.md", filepath.Base(path))
	return !info.IsDir() && err == nil
}

func read(path string) (DailyReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DailyReport{}, err
	}

	t, err := time.Parse("20060102.md", filepath.Base(path))
	if err != nil {
		return DailyReport{}, err
	}

	return unmarshal(data, t)
}

func unmarshal(data []byte, t time.Time) (DailyReport, error) {
	var report DailyReport
	var currentProject string

	projectPattern := regexp.MustCompile(`^- \[.\] (.+)$`)
	taskPattern := regexp.MustCompile(`^\s+- \[(.?)\] ([^\s]+)/([^\s]+) (.+)$`)

	scanner := bufio.NewScanner(bytes.NewReader(data))

LOOP:
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "- 始業 "):
			d, err := parseDuration("- 始業 15:04", line)
			if err != nil {
				return report, err
			}
			report.Attendance.StartedAt = t.Add(d)
		case strings.HasPrefix(line, "- 終業 "):
			d, err := parseDuration("- 終業 15:04", line)
			if err != nil {
				return report, err
			}
			report.Attendance.EndedAt = t.Add(d)
		case strings.HasPrefix(line, "- 休憩 "):
			d, err := parseDuration("- 休憩 15:04", line)
			if err != nil {
				return report, err
			}
			report.Attendance.Break = d
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
			report.Tasks = append(report.Tasks, Task{
				Project:     currentProject,
				Description: matches[4],
				Estimated:   estimate,
				Actual:      actual,
				Completed:   matches[1] == "x",
			})
		case line == "---":
			break LOOP
		}
	}

	report.Attendance.Working = report.Attendance.EndedAt.Sub(report.Attendance.StartedAt) - report.Attendance.Break

	return report, nil
}

func parseDuration(layout string, value string) (time.Duration, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return 0, err
	}
	return t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())), nil
}

func sort(reports []DailyReport) []DailyReport {
	s := slices.Clone(reports)
	slices.SortFunc(s, func(r1, r2 DailyReport) int {
		return r1.Attendance.StartedAt.Compare(r2.Attendance.StartedAt)
	})
	return s
}

func filter(reports []DailyReport, since time.Time, until time.Time) []DailyReport {
	return slices.DeleteFunc(reports, func(report DailyReport) bool {
		return report.Attendance.StartedAt.Before(since) || report.Attendance.StartedAt.After(until)
	})
}
