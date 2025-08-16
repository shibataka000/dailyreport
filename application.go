package main

import (
	"context"
	"time"
)

type Application struct {
	dailyreport *DailyReportRepository
}

func newApplication(dailyreport *DailyReportRepository) *Application {
	return &Application{
		dailyreport: dailyreport,
	}
}

func (app *Application) query(ctx context.Context, since time.Time, until time.Time, query string) ([]byte, error) {
	daily, err := app.dailyreport.list(since, until)
	if err != nil {
		return nil, err
	}
	output := JQOutput{
		Daily:      daily,
		Aggregated: aggregate(daily),
	}
	return jq(ctx, output, query)
}
