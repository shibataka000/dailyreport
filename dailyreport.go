package main

import (
	"time"
)

type DailyReport struct {
	attendance Attendance
	tasks      []Task
}

type Task struct {
	project     string
	description string
	estimate    time.Duration
	actual      time.Duration
	completion  bool
}

type Attendance struct {
	start     time.Time
	end       time.Time
	breakTime time.Duration
}
