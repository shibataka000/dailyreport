# 前提知識

## 日報

ユーザーはその日に行った業務の内容をマークダウン形式のファイルとして記録しています。これを **日報** と呼びます。

日報には以下の情報が記載されています。

- `業務時間` セクションには始業時刻・終業時刻・休憩時間が記載されています。
- `今日のタスク（予定/実績）` セクションにはその日に行ったタスクの名前・所要時間の見積もりと実績・ステータスがプロジェクトごとに To Do リスト形式で記載されています。

以下は日報の例です。

```md
# 日報

## 業務時間

- 始業 09:30
- 終業 17:30
- 休憩 01:00

## 今日のタスク（予定/実績）

- [ ] プロジェクト A
  - [ ] 2.0h/2.5h タスク C
  - [x] 2.0h/1.5h タスク D
- [ ] プロジェクト B
  - [ ] 1.5h/1.25h タスク E
  - [ ] 1.5h/1.75h タスク F

---

# 業務記録
```

- 始業時刻が 9:00、終業時刻が 18:00、休憩時間が 1 時間です。よって労働時間は 7 時間です。
- プロジェクト A のタスク C は、所要時間の見積もりが 2 時間 0 分、実績は 2 時間 30 分でした。ステータスは未完了です。
- プロジェクト A のタスク D は、所要時間の見積もりが 2 時間 0 分、実績は 1 時間 30 分でした。ステータスは完了です。
- プロジェクト B のタスク E は、所要時間の見積もりが 1 時間 30 分、実績は 1 時間 15 分でした。ステータスは未完了です。
- プロジェクト B のタスク F は、所要時間の見積もりが 1 時間 30 分、実績は 1 時間 45 分でした。ステータスは未完了です。

日報はユーザーが指定した任意のディレクトリに `YYYYMMDD.md` という名前で格納されています。 `YYYYMMDD` は日報を作成した年月日です。

# 基本設計

このアプリケーションは日報を集計するための CLI ツールです。

## 入力

このアプリケーションは以下の入力パラメーターを受け取ります。

| パラメーター前 | パラメーター名（短縮形） | 型                                   | デフォルト値               | 説明                               |
| :------------- | :----------------------- | :----------------------------------- | :------------------------- | :--------------------------------- |
| `from`         | N/A                      | 文字列（フォーマットは `YYYYMMDD` ） | アプリケーションの実行日。 | 集計期間の開始日。                 |
| `to`           | N/A                      | 文字列（フォーマットは `YYYYMMDD` ） | アプリケーションの実行日。 | 集計期間の開始日。                 |
| `dir`          | N/A                      | 文字列                               | N/A                        | 日報が格納されているディレクトリ。 |

## 処理

このアプリケーションは以下の処理を行います。

1. ファイルシステムから指定された集計期間の日報を読み込みます。
2. 日報を集計します。どのような集計を行うべきかは後述する出力フォーマットから逆算してください。

## 出力

集計結果を以下のフォーマットに従って標準出力へ出力します。

```md
# サマリー

- 日数 : 3 days
- 業務時間 : 7.00 h
- タスク（予定） : 7.00 h
- タスク（実績） : 7.00 h

# タスク

- [ ] 4.00h / 4.00h プロジェクト A
  - [ ] 2.00h / 2.50h タスク C
  - [x] 2.00h / 1.50h タスク D
- [ ] 3.00h / 3.00h プロジェクト B
  - [ ] 1.50h / 1.25h タスク E
  - [ ] 1.50h / 1.75h タスク F
```

`サマリー` セクションには以下の内容を記載します。

- `日数` には読み取った日報の件数を記載します。
- `業務時間` には読み取った日報の業務時間の合計を記載します。
- `タスク（予定）` には読み取った日報のすべてのタスクの所要時間の見積もりの合計を記載します。
- `タスク（実績）` には読み取った日報のすべてのタスクの所要時間の実績の合計を記載します。

`タスク` セクションには以下の内容を記載します。

- タスクの名前・所要時間の見積もりと実績・ステータスをプロジェクトごとに To Do リスト形式で記載します。プロジェクト名とタスク名が同じタスクが存在する場合、それらは同じタスクとして、所要時間の見積もりと実績を合算して、1 つのタスクに統合します。

## 例

### 例 1

例えば [testdata/case01](../testdata/case01/) の日報に対して以下のコマンドを実行すると仮定します。

```bash
dailyreport --from 20250101 --to 20250107 --dir ./testdata/case01
```

このとき期待される出力は以下のとおりです。

```md
# サマリー

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
```

# 詳細設計

## 全体

- このアプリケーションは Go 言語で実装します。
- このアプリケーションは "ApplicationService" , "Infrastructure" , "Domain" の 3 層で構成される Layered Architecture として実装します。

## `main` package

このアプリケーションのエンドポイントです。

- ソースコードは [cmd/dailyreport](../cmd/dailyreport) に配置します。
- CLI ツールのフレームワークとして https://github.com/spf13/cobra を使用します。

## `github.com/shibataka000/dailyreport` package

日報に関するパッケージです。

- ソースコードは [/](../) に配置します。
- 以下の構造体・関数・メソッドを含みます。

```go
// ApplicationService 層のオブジェクトです。
// 主な役割はパッケージの利用者に対してこのパッケージの機能を提供することです。
// 構造体のフィールドの定義は省略します。実装時に必要なフィールドを追加してください。
type ApplicationService struct { /*...*/ }

// Infrastructure 層のオブジェクトです。
// 主な役割はファイルシステムから日報を読み取ることです。
// 構造体のフィールドの定義は省略します。実装時に必要なフィールドを追加してください。
type ReportRepository struct { /*...*/ }

// Domain 層のオブジェクトです。
// 日報を表すドメインオブジェクトです。
// 構造体のフィールドの定義は省略します。実装時に必要なフィールドを追加してください。
type Report struct { /*...*/ }

// 新しい ApplicationService オブジェクトを作成します。
func NewApplicationService(ctx context.Context, repository *ReportRepository) (*ApplicationService, error)

// 新しい ReportRepository オブジェクトを作成します。
func NewReportRepository(ctx context.Context, dir string) (*ReportRepository, error)

// 指定された期間の日報を読み取り、その集計結果を所定のフォーマットの文字列として返します。
func (a *ApplicationService) Aggregate(ctx context.Context, from time.Time, end time.Time) (string, error)

// 指定された日付の日報を読み取ります。
func (r *ReportRepository) Read(ctx context.Context, date time.Time) (Report, error)
```
