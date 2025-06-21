package dailyreport_test

import (
	"context"
	"testing"
	"time"

	"github.com/shibataka000/dailyreport"
)

func TestApplicationService_Aggregate_Case01(t *testing.T) {
	ctx := context.Background()
	repo, err := dailyreport.NewReportRepository(ctx, "testdata/case01")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	app, err := dailyreport.NewApplicationService(ctx, repo)
	if err != nil {
		t.Fatalf("failed to create application service: %v", err)
	}
	from, _ := time.Parse("20060102", "20250101")
	to, _ := time.Parse("20060102", "20250107")
	got, err := app.Aggregate(ctx, from, to)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	want := `# サマリー

- 日数 : 2 days
- 業務時間 : 14.00 h
- タスク（予定） : 11.00 h
- タスク（実績） : 14.00 h

# タスク

- [ ] 8.00h / 8.00h プロジェクト A
  - [ ] 4.00h / 5.00h タスク C
  - [x] 4.00h / 3.00h タスク D
- [ ] 3.00h / 3.00h プロジェクト B1
  - [ ] 1.50h / 1.25h タスク E
  - [ ] 1.50h / 1.75h タスク F
- [ ] 0.00h / 3.00h プロジェクト B2
  - [ ] 0.00h / 1.25h タスク E
  - [ ] 0.00h / 1.75h タスク F
`
	if got != want {
		t.Errorf("Aggregate output mismatch.\nGot:\n%s\nWant:\n%s", got, want)
	}
}
