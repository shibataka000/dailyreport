package main

import "time"

// JQOutput represents the output structure containing a list of daily reports
// and an aggregated report. The 'daily' field holds multiple DailyReport entries,
// while the 'aggregated' field summarizes the data in an AggregatedReport.
type JQOutput struct {
	Daily      []DailyReport    `json:"daily"`
	Aggregated AggregatedReport `json:"aggregated"`
}

// DailyReport represents a daily report containing attendance information and a list of tasks.
type DailyReport struct {
	Attendance Attendance `json:"attendance"`
	Tasks      []Task     `json:"tasks"`
}

// AggregatedReport represents a collection of tasks that have been aggregated for reporting purposes.
type AggregatedReport struct {
	Tasks []Task `json:"tasks"`
}

// Attendance represents a record of an individual's attendance for a given period.
// It includes the start and end times, the duration of breaks taken, and the total working duration.
type Attendance struct {
	StartedAt time.Time     `json:"started_at"`
	EndedAt   time.Time     `json:"ended_at"`
	Break     time.Duration `json:"break"`
	Working   time.Duration `json:"working"`
}

// Task represents a work item with associated project, description, estimated and actual durations,
// and completion status. It is used to track progress and time spent on specific tasks.
type Task struct {
	Project     string        `json:"project"`
	Description string        `json:"description"`
	Estimated   time.Duration `json:"estimated"`
	Actual      time.Duration `json:"actual"`
	Completed   bool          `json:"completed"`
}
