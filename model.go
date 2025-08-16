package main

import "time"

type JQOutput struct {
	Daily      []DailyReport    `json:"daily"`
	Aggregated AggregatedReport `json:"aggregated"`
}

type DailyReport struct {
	Attendance Attendance `json:"attendance"`
	Tasks      []Task     `json:"tasks"`
}

type AggregatedReport struct {
	Tasks []Task `json:"tasks"`
}

type Attendance struct {
	StartedAt time.Time     `json:"started_at"`
	EndedAt   time.Time     `json:"ended_at"`
	Break     time.Duration `json:"break"`
	Working   time.Duration `json:"working"`
}

type Task struct {
	Project     string        `json:"project"`
	Description string        `json:"description"`
	Estimated   time.Duration `json:"estimated"`
	Actual      time.Duration `json:"actual"`
	Completed   bool          `json:"completed"`
}
