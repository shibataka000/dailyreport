# Daily report

[![Test](https://github.com/shibataka000/dailyreport/actions/workflows/test.yaml/badge.svg)](https://github.com/shibataka000/dailyreport/actions/workflows/test.yaml)

CLI tool to control daily report.

## Usage

### Show work time

```
Show work time in daily report.

Usage:
  dailyreport show worktime [flags]

Flags:
      --dir string        Directory where daily report file exists. [$DR_DIR]
  -e, --end-at string     End of daily report date range. [$DR_END_AT]
  -h, --help              help for worktime
      --project string    Show only tasks which project name is this. [$DR_PROJECT]
  -s, --start-at string   Start of daily report date range. [$DR_START_AT]
```

### Show tasks

```
Show tasks in daily report.

Usage:
  dailyreport show tasks [flags]

Flags:
      --dir string        Directory where daily report file exists. [$DR_DIR]
  -e, --end-at string     End of daily report date range. [$DR_END_AT]
  -h, --help              help for tasks
      --project string    Show only tasks which project name is this. [$DR_PROJECT]
  -s, --start-at string   Start of daily report date range. [$DR_START_AT]
```

## Install

```
go install github.com/shibataka000/dailyreport@master
```
