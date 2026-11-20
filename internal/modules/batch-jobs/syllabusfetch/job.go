package syllabusfetch

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Summary はジョブ 1 回分の処理件数。
type Summary struct {
	Listed    int
	Inserted  int
	Updated   int
	Unchanged int
	Failed    int
}

// detailRetryWait は詳細ページの取得が一時的に失敗したとき、やり直すまでに待つ時間。
const detailRetryWait = 10 * time.Second

type Job struct {
	client    *Client
	store     *Store
	now       func() time.Time
	retryWait time.Duration
}

func NewJob(client *Client, store *Store) *Job {
	return &Job{client: client, store: store, now: time.Now, retryWait: detailRetryWait}
}

// Run は指定年度のシラバスを全件取得して保存する。
// 1 件の詳細ページの取得・保存に失敗しても残りの科目は処理を続け、失敗件数を Summary に含める。
func (j *Job) Run(ctx context.Context, year int) (Summary, error) {
	ids, err := j.client.ListLessonIDs(ctx, year)
	if err != nil {
		return Summary{}, fmt.Errorf("list lesson ids: %w", err)
	}
	if len(ids) == 0 {
		return Summary{}, fmt.Errorf("no syllabus found for year %d", year)
	}

	summary := Summary{Listed: len(ids)}
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return summary, err
		}
		result, err := j.fetchAndSave(ctx, id, year)
		if err != nil {
			summary.Failed++
			slog.ErrorContext(ctx, "failed to fetch syllabus", "lesson_id", id, "year", year, "error", err)
			continue
		}
		switch result {
		case SaveInserted:
			summary.Inserted++
		case SaveUpdated:
			summary.Updated++
		case SaveUnchanged:
			summary.Unchanged++
		}
	}
	return summary, nil
}

func (j *Job) fetchAndSave(ctx context.Context, lessonID string, year int) (SaveResult, error) {
	body, err := j.fetchDetailWithRetry(ctx, lessonID, year)
	if err != nil {
		return 0, err
	}
	syllabus, unknown, err := ParseDetailPage(bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	if len(unknown) > 0 {
		slog.WarnContext(ctx, "syllabus has unknown items", "lesson_id", lessonID, "year", year, "labels", unknown)
	}
	// エラーページなどを空のシラバスとして保存しないため、授業名が取れない場合は失敗として扱う
	if syllabus.Name == "" {
		return 0, fmt.Errorf("lesson name not found in detail page")
	}
	syllabus.LessonID = lessonID
	syllabus.Year = year
	return j.store.Save(ctx, syllabus, j.now())
}

// fetchDetailWithRetry は詳細ページを取得する。5xx やタイムアウトのときだけ、間隔を空けて 1 回やり直す。
// 一時的な失敗 1 件でジョブ全体が失敗扱いになり、毎週アラートが鳴るのを避けるため。
func (j *Job) fetchDetailWithRetry(ctx context.Context, lessonID string, year int) ([]byte, error) {
	body, err := j.client.FetchDetail(ctx, lessonID, year)
	if err == nil || !isRetryable(err) {
		return body, err
	}
	slog.WarnContext(ctx, "retry fetching syllabus", "lesson_id", lessonID, "year", year, "error", err)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(j.retryWait):
	}
	return j.client.FetchDetail(ctx, lessonID, year)
}
