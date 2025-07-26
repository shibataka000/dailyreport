package main

import (
	"time"
)

type DailyReport struct {
	attendance Attendance
	tasks      []Task
}

type Attendance struct {
	start     time.Time
	end       time.Time
	breakTime time.Duration
}

type Task struct {
	project     string
	description string
	estimate    time.Duration
	actual      time.Duration
	completion  bool
}
