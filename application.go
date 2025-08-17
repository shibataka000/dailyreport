package main

import (
	"context"
	"time"
)

// Application encapsulates the dependencies required for the daily report application.
// It holds a reference to DailyReportRepository for managing daily report data.
type Application struct {
	dailyreport *DailyReportRepository
}

func newApplication(dailyreport *DailyReportRepository) *Application {
	return &Application{
		dailyreport: dailyreport,
	}
}

func (app *Application) query(ctx context.Context, since time.Time, until time.Time, query string) (string, error) {
	daily, err := app.dailyreport.list(since, until)
	if err != nil {
		return "", err
	}
	output := JQOutput{
		Daily:      daily,
		Aggregated: aggregate(daily),
	}
	return jq(ctx, output, query)
}
