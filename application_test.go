package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		name   string
		dir    string
		since  time.Time
		until  time.Time
		query  string
		output []byte
	}{
		{
			name:  "Case01",
			dir:   "./testdata",
			since: time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			until: time.Date(2025, 1, 3, 9, 15, 0, 0, time.UTC),
			query: ".",
			output: []byte(`{
  "daily": [
    {
      "attendance": {
        "started_at": "2025-01-01T09:30:00Z",
        "ended_at": "2025-01-01T17:30:00Z",
        "break": 3600000000000,
        "working": 25200000000000
      },
      "tasks": [
        {
          "project": "プロジェクト A",
          "description": "タスク C",
          "estimated": 7200000000000,
          "actual": 9000000000000,
          "completed": false
        },
        {
          "project": "プロジェクト A",
          "description": "タスク D",
          "estimated": 7200000000000,
          "actual": 5400000000000,
          "completed": false
        },
        {
          "project": "プロジェクト B1",
          "description": "タスク E",
          "estimated": 5400000000000,
          "actual": 4500000000000,
          "completed": false
        },
        {
          "project": "プロジェクト B1",
          "description": "タスク F",
          "estimated": 5400000000000,
          "actual": 6300000000000,
          "completed": false
        }
      ]
    },
    {
      "attendance": {
        "started_at": "2025-01-03T09:15:00Z",
        "ended_at": "2025-01-03T17:45:00Z",
        "break": 5400000000000,
        "working": 25200000000000
      },
      "tasks": [
        {
          "project": "プロジェクト A",
          "description": "タスク C",
          "estimated": 7200000000000,
          "actual": 9000000000000,
          "completed": false
        },
        {
          "project": "プロジェクト A",
          "description": "タスク D",
          "estimated": 7200000000000,
          "actual": 5400000000000,
          "completed": true
        },
        {
          "project": "プロジェクト B2",
          "description": "タスク E",
          "estimated": 0,
          "actual": 4500000000000,
          "completed": false
        },
        {
          "project": "プロジェクト B2",
          "description": "タスク F",
          "estimated": 0,
          "actual": 6300000000000,
          "completed": false
        }
      ]
    }
  ],
  "aggregated": {
    "tasks": [
      {
        "project": "プロジェクト A",
        "description": "タスク C",
        "estimated": 14400000000000,
        "actual": 18000000000000,
        "completed": false
      },
      {
        "project": "プロジェクト A",
        "description": "タスク D",
        "estimated": 14400000000000,
        "actual": 10800000000000,
        "completed": true
      },
      {
        "project": "プロジェクト B1",
        "description": "タスク E",
        "estimated": 5400000000000,
        "actual": 4500000000000,
        "completed": false
      },
      {
        "project": "プロジェクト B1",
        "description": "タスク F",
        "estimated": 5400000000000,
        "actual": 6300000000000,
        "completed": false
      },
      {
        "project": "プロジェクト B2",
        "description": "タスク E",
        "estimated": 0,
        "actual": 4500000000000,
        "completed": false
      },
      {
        "project": "プロジェクト B2",
        "description": "タスク F",
        "estimated": 0,
        "actual": 6300000000000,
        "completed": false
      }
    ]
  }
}
`),
		},
		{
			name:  "Case02",
			dir:   "./testdata",
			since: time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			until: time.Date(2025, 1, 3, 9, 14, 0, 0, time.UTC),
			query: ".",
			output: []byte(`{
  "daily": [
    {
      "attendance": {
        "started_at": "2025-01-01T09:30:00Z",
        "ended_at": "2025-01-01T17:30:00Z",
        "break": 3600000000000,
        "working": 25200000000000
      },
      "tasks": [
        {
          "project": "プロジェクト A",
          "description": "タスク C",
          "estimated": 7200000000000,
          "actual": 9000000000000,
          "completed": false
        },
        {
          "project": "プロジェクト A",
          "description": "タスク D",
          "estimated": 7200000000000,
          "actual": 5400000000000,
          "completed": false
        },
        {
          "project": "プロジェクト B1",
          "description": "タスク E",
          "estimated": 5400000000000,
          "actual": 4500000000000,
          "completed": false
        },
        {
          "project": "プロジェクト B1",
          "description": "タスク F",
          "estimated": 5400000000000,
          "actual": 6300000000000,
          "completed": false
        }
      ]
    }
  ],
  "aggregated": {
    "tasks": [
      {
        "project": "プロジェクト A",
        "description": "タスク C",
        "estimated": 7200000000000,
        "actual": 9000000000000,
        "completed": false
      },
      {
        "project": "プロジェクト A",
        "description": "タスク D",
        "estimated": 7200000000000,
        "actual": 5400000000000,
        "completed": false
      },
      {
        "project": "プロジェクト B1",
        "description": "タスク E",
        "estimated": 5400000000000,
        "actual": 4500000000000,
        "completed": false
      },
      {
        "project": "プロジェクト B1",
        "description": "タスク F",
        "estimated": 5400000000000,
        "actual": 6300000000000,
        "completed": false
      }
    ]
  }
}
`),
		},
		{
			name:  "Case03",
			dir:   "./testdata",
			since: time.Date(2025, 1, 1, 9, 31, 0, 0, time.UTC),
			until: time.Date(2025, 1, 3, 9, 15, 0, 0, time.UTC),
			query: ".",
			output: []byte(`{
  "daily": [
    {
      "attendance": {
        "started_at": "2025-01-03T09:15:00Z",
        "ended_at": "2025-01-03T17:45:00Z",
        "break": 5400000000000,
        "working": 25200000000000
      },
      "tasks": [
        {
          "project": "プロジェクト A",
          "description": "タスク C",
          "estimated": 7200000000000,
          "actual": 9000000000000,
          "completed": false
        },
        {
          "project": "プロジェクト A",
          "description": "タスク D",
          "estimated": 7200000000000,
          "actual": 5400000000000,
          "completed": true
        },
        {
          "project": "プロジェクト B2",
          "description": "タスク E",
          "estimated": 0,
          "actual": 4500000000000,
          "completed": false
        },
        {
          "project": "プロジェクト B2",
          "description": "タスク F",
          "estimated": 0,
          "actual": 6300000000000,
          "completed": false
        }
      ]
    }
  ],
  "aggregated": {
    "tasks": [
      {
        "project": "プロジェクト A",
        "description": "タスク C",
        "estimated": 7200000000000,
        "actual": 9000000000000,
        "completed": false
      },
      {
        "project": "プロジェクト A",
        "description": "タスク D",
        "estimated": 7200000000000,
        "actual": 5400000000000,
        "completed": true
      },
      {
        "project": "プロジェクト B2",
        "description": "タスク E",
        "estimated": 0,
        "actual": 4500000000000,
        "completed": false
      },
      {
        "project": "プロジェクト B2",
        "description": "タスク F",
        "estimated": 0,
        "actual": 6300000000000,
        "completed": false
      }
    ]
  }
}
`),
		},
		{
			name:  "Case04",
			dir:   "./testdata",
			since: time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			until: time.Date(2025, 1, 3, 9, 15, 0, 0, time.UTC),
			query: ".daily[0].attendance.working",
			output: []byte(`25200000000000
`),
		},
		{
			name:  "Case05",
			dir:   "./testdata",
			since: time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			until: time.Date(2025, 1, 3, 9, 15, 0, 0, time.UTC),
			query: "\"x\"",
			output: []byte(`x
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			app := newApplication(newDailyReportRepository(tt.dir))
			output, err := app.query(context.Background(), tt.since, tt.until, tt.query)
			require.NoError(err)
			require.Equal(string(tt.output), string(output))
		})
	}
}
