# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Install

```bash
# Build the binary
make build

# Install the binary to your GOPATH
make install
```

### Testing

```bash
# Run all tests
make test

# Run a specific test
go test -run TestQuery/Case01

# Clean test cache
make clean
```

### Linting

```bash
# Run linter
make lint
```

## Code Architecture

The dailyreport tool is a CLI application that processes and queries daily report markdown files. It follows a layered architecture:

1. **Main Layer** (`main.go`): Entry point that sets up the CLI using Cobra, processes command-line arguments, and starts the application.

2. **Application Layer** (`application.go`): Contains the core application logic, orchestrating the use cases by connecting repositories and services.

3. **Domain Layer** (`model.go`): Defines the domain models including:

   - `DailyReport`: Contains attendance information and tasks
   - `Task`: Represents work items with project, description, estimated/actual durations
   - `Attendance`: Tracks work times and breaks
   - `JQOutput`: Output structure for queries containing daily reports and aggregated reports

4. **Service Layer** (`service.go`): Provides utility functions like:

   - `aggregate()`: Aggregates tasks from multiple reports
   - `jq()`: Executes jq queries on the report data

5. **Infrastructure Layer** (`infrastructure.go`): Handles data access and persistence:
   - `DailyReportRepository`: Manages reading daily reports from the filesystem
   - Parses markdown files with a specific format (YYYYMMDD.md) into structured data

### Daily Report Format

Reports are stored as markdown files with a specific format (YYYYMMDD.md):

```md
# 日報

## 業務時間

- 始業 HH:MM
- 終業 HH:MM
- 休憩 HH:MM

## 今日のタスク（予定/実績）

- [ ] プロジェクト名 a
  - [ ] 予定時間/実績時間 タスク説明
  - [x] 予定時間/実績時間 タスク説明 (完了済み)

---

# 業務記録
```

### Command Line Usage

The tool can be run with the following flags:

```
dailyreport --dir=<directory> --since=<date> --until=<date> --query=<jq query>
```

- `--dir`: Directory containing daily reports (default: current directory)
- `--since`: Start date (inclusive) for querying reports (format: YYYY-MM-DD)
- `--until`: End date (inclusive) for querying reports (format: YYYY-MM-DD)
- `--query`: JQ query string(s) to filter reports (can specify multiple)

The tool reads daily report markdown files from the specified directory, parses them into structured data, applies JQ queries, and outputs the results.
