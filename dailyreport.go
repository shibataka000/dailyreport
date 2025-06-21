// Package dailyreport provides domain logic for daily report aggregation.
package dailyreport

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ApplicationService provides application-level services for daily report aggregation.
type ApplicationService struct {
	repository *ReportRepository
}

// ReportRepository handles reading daily reports from the file system.
type ReportRepository struct {
	dir string
}

// Report represents a daily report domain object.
type Report struct {
	Date      time.Time
	Start     time.Time // 始業時刻
	End       time.Time // 終業時刻
	Break     float64   // 休憩時間（h）
	WorkHours float64   // 労働時間（h）
	Tasks     []Task
}

// Task represents a task in a daily report.
type Task struct {
	Project   string
	Name      string
	Estimate  float64 // 見積もり（h）
	Actual    float64 // 実績（h）
	Completed bool
}

// NewApplicationService creates a new ApplicationService object.
func NewApplicationService(_ context.Context, repository *ReportRepository) (*ApplicationService, error) {
	return &ApplicationService{repository: repository}, nil
}

// NewReportRepository creates a new ReportRepository object.
func NewReportRepository(_ context.Context, dir string) (*ReportRepository, error) {
	return &ReportRepository{dir: dir}, nil
}

// Aggregate reads daily reports in the specified period and returns the aggregation result as a formatted string.
func (a *ApplicationService) Aggregate(ctx context.Context, from time.Time, to time.Time) (string, error) {
	reports := []Report{}
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		report, err := a.repository.Read(ctx, d)
		if err != nil {
			if strings.Contains(err.Error(), "file not found") {
				continue
			}
			return "", err
		}
		reports = append(reports, report)
	}
	return aggregateReportsToString(reports), nil
}

// Read reads a daily report for the specified date.
func (r *ReportRepository) Read(_ context.Context, date time.Time) (Report, error) {
	filename := filepath.Join(r.dir, date.Format("20060102")+".md")
	b, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return Report{}, fmt.Errorf("file not found: %s", filename)
		}
		return Report{}, err
	}
	return parseReportMarkdown(date, string(b))
}

// --- Domain 層 ---

// parseReportMarkdown parses a markdown file content into a Report.
func parseReportMarkdown(date time.Time, content string) (Report, error) {
	dr := Report{Date: date}
	lines := strings.Split(content, "\n")

	var (
		startRe = regexp.MustCompile(`始業\s+(\d{2}):(\d{2})`)
		endRe   = regexp.MustCompile(`終業\s+(\d{2}):(\d{2})`)
		breakRe = regexp.MustCompile(`休憩\s+(\d{2}):(\d{2})|休憩\s+(\d+\.\d+|\d+)`)
		projRe  = regexp.MustCompile(`^[-*] \[.\] ([^0-9].*)$`)
		taskRe  = regexp.MustCompile(`^\s*- \[( |x)\] ([0-9.]+)h/([0-9.]+)h (.+)$`)
	)

	var currentProject string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if m := startRe.FindStringSubmatch(line); m != nil {
			hour, _ := strconv.Atoi(m[1])
			minute, _ := strconv.Atoi(m[2])
			dr.Start = time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.Local)
		}
		if m := endRe.FindStringSubmatch(line); m != nil {
			hour, _ := strconv.Atoi(m[1])
			minute, _ := strconv.Atoi(m[2])
			dr.End = time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.Local)
		}
		if m := breakRe.FindStringSubmatch(line); m != nil {
			if m[1] != "" && m[2] != "" {
				h, _ := strconv.Atoi(m[1])
				minute, _ := strconv.Atoi(m[2])
				dr.Break = float64(h) + float64(minute)/60.0
			} else if m[3] != "" {
				f, _ := strconv.ParseFloat(m[3], 64)
				dr.Break = f
			}
		}
		if m := projRe.FindStringSubmatch(line); m != nil {
			currentProject = m[1]
		}
		if m := taskRe.FindStringSubmatch(line); m != nil {
			estimate, _ := strconv.ParseFloat(m[2], 64)
			actual, _ := strconv.ParseFloat(m[3], 64)
			dr.Tasks = append(dr.Tasks, Task{
				Project:   strings.TrimSpace(currentProject),
				Name:      m[4],
				Estimate:  estimate,
				Actual:    actual,
				Completed: m[1] == "x",
			})
		}
	}
	if !dr.Start.IsZero() && !dr.End.IsZero() {
		dur := dr.End.Sub(dr.Start).Hours() - dr.Break
		if dur < 0 {
			dur = 0
		}
		dr.WorkHours = dur
	}
	return dr, nil
}

// aggregateReportsToString converts aggregated reports to the specified output format.
func aggregateReportsToString(reports []Report) string {
	if len(reports) == 0 {
		return "# サマリー\n\n- 日数 : 0 days\n- 業務時間 : 0.00 h\n- タスク（予定） : 0.00 h\n- タスク（実績） : 0.00 h\n\n# タスク\n"
	}

	totalDays := len(reports)
	totalWork := 0.0
	totalEstimate := 0.0
	totalActual := 0.0
	// プロジェクト単位でタスクをまとめる構造体を用意
	type projectSummary struct {
		Estimate float64
		Actual   float64
		Tasks    map[string]*Task
		Order    []string // タスク名の出現順
	}
	projectMap := map[string]*projectSummary{}
	projectOrder := []string{}

	for _, dr := range reports {
		totalWork += dr.WorkHours
		for _, t := range dr.Tasks {
			totalEstimate += t.Estimate
			totalActual += t.Actual
			if _, ok := projectMap[t.Project]; !ok {
				projectMap[t.Project] = &projectSummary{Tasks: map[string]*Task{}}
				projectOrder = append(projectOrder, t.Project)
			}
			ps := projectMap[t.Project]
			if _, ok := ps.Tasks[t.Name]; !ok {
				ps.Tasks[t.Name] = &Task{
					Project:   t.Project,
					Name:      t.Name,
					Estimate:  t.Estimate,
					Actual:    t.Actual,
					Completed: t.Completed,
				}
				ps.Order = append(ps.Order, t.Name)
			} else {
				task := ps.Tasks[t.Name]
				task.Estimate += t.Estimate
				task.Actual += t.Actual
				if t.Completed {
					task.Completed = true
				}
			}
			ps.Estimate += t.Estimate
			ps.Actual += t.Actual
		}
	}

	var sb strings.Builder
	sb.WriteString("# サマリー\n\n")
	sb.WriteString(fmt.Sprintf("- 日数 : %d days\n", totalDays))
	sb.WriteString(fmt.Sprintf("- 業務時間 : %.2f h\n", totalWork))
	sb.WriteString(fmt.Sprintf("- タスク（予定） : %.2f h\n", totalEstimate))
	sb.WriteString(fmt.Sprintf("- タスク（実績） : %.2f h\n\n", totalActual))

	sb.WriteString("# タスク\n\n")
	for _, project := range projectOrder {
		ps := projectMap[project]
		projDone := true
		for _, t := range ps.Tasks {
			if !t.Completed {
				projDone = false
				break
			}
		}
		if projDone {
			sb.WriteString(fmt.Sprintf("- [x] %.2fh / %.2fh %s\n", ps.Estimate, ps.Actual, project))
		} else {
			sb.WriteString(fmt.Sprintf("- [ ] %.2fh / %.2fh %s\n", ps.Estimate, ps.Actual, project))
		}
		for _, name := range ps.Order {
			t := ps.Tasks[name]
			box := "[ ]"
			if t.Completed {
				box = "[x]"
			}
			sb.WriteString(fmt.Sprintf("  - %s %.2fh / %.2fh %s\n", box, t.Estimate, t.Actual, t.Name))
		}
	}
	return sb.String()
}
