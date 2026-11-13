package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/fun-dotto/server/internal/shared/syllabussource"
	"gorm.io/gorm"
)

// SyllabusSource は教務システムからシラバスを取得する。
type SyllabusSource interface {
	ListAll(ctx context.Context, year int) ([]syllabussource.ListEntry, error)
	Detail(ctx context.Context, lessonID string, year int) (syllabussource.Detail, error)
}

// Notifier は差分の通知を送る。
type Notifier interface {
	Post(ctx context.Context, text string) error
}

// SyllabusChange は既存の科目が古いバージョンに紐付いている授業の差分。
type SyllabusChange struct {
	LessonID      string
	Year          int
	Name          string
	SubjectID     string
	ChangedFields []string
}

// SyllabusFetchSummary はシラバス取得ジョブの結果。
type SyllabusFetchSummary struct {
	Listed    int
	Inserted  int
	Unchanged int
	Failed    int
	Changes   []SyllabusChange
}

type SyllabusFetchService struct {
	db       *gorm.DB
	source   SyllabusSource
	notifier Notifier
	now      func() time.Time
}

func NewSyllabusFetchService(db *gorm.DB, source SyllabusSource, notifier Notifier) *SyllabusFetchService {
	return &SyllabusFetchService{db: db, source: source, notifier: notifier, now: time.Now}
}

// Fetch は year の一覧・詳細を取得して保存する。
// 内容が既存のバージョンと同じなら fetched_at のみ更新し、異なれば新しいバージョンとして追加する。
// 新しいバージョンを追加した授業のうち、既存の科目が古いバージョンに紐付いているものを通知する。
func (s *SyllabusFetchService) Fetch(ctx context.Context, year int) (SyllabusFetchSummary, error) {
	var summary SyllabusFetchSummary
	entries, err := s.source.ListAll(ctx, year)
	if err != nil {
		return summary, fmt.Errorf("list syllabuses: %w", err)
	}
	summary.Listed = len(entries)
	slog.InfoContext(ctx, "listed syllabuses", "year", year, "count", len(entries))

	for _, e := range entries {
		detail, err := s.source.Detail(ctx, e.LessonID, year)
		if err != nil {
			summary.Failed++
			slog.ErrorContext(ctx, "failed to fetch syllabus detail", "lesson_id", e.LessonID, "year", year, "error", err)
			continue
		}
		m, unknown := syllabussource.ToModel(year, e, detail)
		if len(unknown) > 0 {
			slog.WarnContext(ctx, "unknown syllabus labels", "lesson_id", e.LessonID, "year", year, "labels", unknown)
		}
		change, inserted, err := s.save(ctx, &m)
		if err != nil {
			summary.Failed++
			slog.ErrorContext(ctx, "failed to save syllabus", "lesson_id", e.LessonID, "year", year, "error", err)
			continue
		}
		if inserted {
			summary.Inserted++
		} else {
			summary.Unchanged++
		}
		if change != nil {
			summary.Changes = append(summary.Changes, *change)
		}
	}

	if len(summary.Changes) > 0 && s.notifier != nil {
		for _, text := range formatSyllabusChanges(year, summary.Changes) {
			if err := s.notifier.Post(ctx, text); err != nil {
				return summary, fmt.Errorf("notify syllabus changes: %w", err)
			}
		}
	}
	if summary.Failed > 0 {
		return summary, fmt.Errorf("%d syllabuses failed", summary.Failed)
	}
	return summary, nil
}

// save は m を保存し、新しいバージョンを追加したかと、古いバージョンに紐付いた科目の差分を返す。
func (s *SyllabusFetchService) save(ctx context.Context, m *model.Syllabus) (*SyllabusChange, bool, error) {
	m.ContentHash = m.ComputeContentHash()
	m.FetchedAt = s.now()

	var change *SyllabusChange
	inserted := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Syllabus{}).
			Where("lesson_id = ? AND year = ? AND content_hash = ?", m.LessonID, m.Year, m.ContentHash).
			Update("fetched_at", m.FetchedAt)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			return nil
		}
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		inserted = true

		// 同じ授業の古いバージョンに紐付いている科目
		var subjects []model.Subject
		err := tx.Preload("Syllabus").
			Joins("JOIN syllabuses ON syllabuses.id = subjects.syllabus_id").
			Where("syllabuses.lesson_id = ? AND syllabuses.year = ? AND syllabuses.id <> ?", m.LessonID, m.Year, m.ID.String()).
			Limit(1).Find(&subjects).Error
		if err != nil {
			return err
		}
		if len(subjects) == 0 {
			return nil
		}
		subject := subjects[0]
		var fields []string
		if subject.Syllabus != nil {
			for _, f := range model.ChangedSyllabusFields(subject.Syllabus, m) {
				fields = append(fields, f.Label)
			}
		}
		change = &SyllabusChange{
			LessonID:      m.LessonID,
			Year:          m.Year,
			Name:          m.Name,
			SubjectID:     subject.ID.String(),
			ChangedFields: fields,
		}
		return nil
	})
	return change, inserted, err
}

// changesPerMessage は 1 メッセージに含める差分の件数。
const changesPerMessage = 30

// formatSyllabusChanges は通知の本文を changesPerMessage 件ずつに分けて返す。
func formatSyllabusChanges(year int, changes []SyllabusChange) []string {
	var texts []string
	for start := 0; start < len(changes); start += changesPerMessage {
		end := min(start+changesPerMessage, len(changes))
		var b strings.Builder
		if start == 0 {
			fmt.Fprintf(&b, "*%d年度のシラバスに変更があります（%d件）*\n", year, len(changes))
			b.WriteString("科目は古いバージョンに紐付いたままです。admin-api で新しいバージョンに付け替えてください。\n")
		}
		for _, c := range changes[start:end] {
			fields := strings.Join(c.ChangedFields, "、")
			if fields == "" {
				fields = "（なし）"
			}
			fmt.Fprintf(&b, "• `%s` %d %s\n    変更項目: %s\n    科目ID: `%s`\n", c.LessonID, c.Year, c.Name, fields, c.SubjectID)
		}
		texts = append(texts, b.String())
	}
	return texts
}
