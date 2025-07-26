package main

import (
	"os"
	"path/filepath"
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

	date, err := time.Parse("20060102.md", filepath.Base(path))
	if err != nil {
		return DailyReport{}, err
	}

	return unmarshal(date, data)
}

func unmarshal(date time.Time, data []byte) (DailyReport, error) {
	return DailyReport{}, nil
}
