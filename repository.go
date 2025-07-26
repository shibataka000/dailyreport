package main

import "os"

type DailyReportRepository struct{}

func newDailyReportRepository() (*DailyReportRepository, error) {
	return &DailyReportRepository{}, nil
}

func (r *DailyReportRepository) read(path string) (DailyReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DailyReport{}, err
	}
	return unmarshal(data)
}
