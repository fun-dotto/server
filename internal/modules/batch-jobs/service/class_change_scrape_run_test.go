package service

import (
	"context"
	"testing"
	"time"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

const scrapeTestHTML = `
<table>
  <tr>
    <td data-col-responsive-title="日付">05/12</td>
    <td data-col-responsive-title="曜日">火</td>
    <td data-col-responsive-title="時限">3時限</td>
    <td data-col-responsive-title="授業名">情報数学</td>
    <td data-col-responsive-title="キャンパス">本校</td>
    <td data-col-responsive-title="代表教職員">教員A</td>
    <td data-col-responsive-title="休講コメント">補講あり 後日連絡</td>
  </tr>
  <tr>
    <td data-col-responsive-title="日付">05/13</td>
    <td data-col-responsive-title="曜日">水</td>
    <td data-col-responsive-title="時限">1時限</td>
    <td data-col-responsive-title="授業名">未知の科目</td>
    <td data-col-responsive-title="キャンパス">本校</td>
    <td data-col-responsive-title="代表教職員">教員B</td>
    <td data-col-responsive-title="休講コメント">補講なし</td>
  </tr>
  <tr>
    <td data-col-responsive-title="日付">05/20</td>
    <td data-col-responsive-title="曜日">水</td>
    <td data-col-responsive-title="時限">2時限</td>
    <td data-col-responsive-title="授業名">情報数学</td>
    <td data-col-responsive-title="キャンパス">本校</td>
    <td data-col-responsive-title="代表教職員">教員A</td>
    <td data-col-responsive-title="補講コメント">5/12 の補講</td>
    <td data-col-responsive-title="教室名">R791</td>
  </tr>
  <tr>
    <td data-col-responsive-title="日付">05/21</td>
    <td data-col-responsive-title="曜日">木</td>
    <td data-col-responsive-title="時限">4時限</td>
    <td data-col-responsive-title="授業名">情報数学</td>
    <td data-col-responsive-title="キャンパス">本校</td>
    <td data-col-responsive-title="代表教職員">教員A</td>
    <td data-col-responsive-title="移動元">R791</td>
    <td data-col-responsive-title="移動先">R792</td>
  </tr>
</table>
`

type stubScraperClient struct {
	html string
}

func (c *stubScraperClient) FetchClassChangeHTML(_ context.Context, _, _ string) (string, error) {
	return c.html, nil
}

type stubSubjectRepository struct {
	subjects []domain.SubjectRef
}

func (r *stubSubjectRepository) ListSyllabusNames(_ context.Context) ([]domain.SubjectRef, error) {
	return r.subjects, nil
}

type stubRoomRepository struct {
	rooms []domain.RoomRef
}

func (r *stubRoomRepository) ListRooms(_ context.Context) ([]domain.RoomRef, error) {
	return r.rooms, nil
}

type stubCancelledClassRepository struct {
	existing map[string]struct{}
	inserted []academicdomain.CancelledClass
}

func cancelledKey(subjectID string, date time.Time, period string) string {
	return subjectID + "|" + date.Format("2006-01-02") + "|" + period
}

func (r *stubCancelledClassRepository) Exists(_ context.Context, subjectID string, date time.Time, period string) (bool, error) {
	_, ok := r.existing[cancelledKey(subjectID, date, period)]
	return ok, nil
}

func (r *stubCancelledClassRepository) Insert(_ context.Context, cc academicdomain.CancelledClass) (academicdomain.CancelledClass, error) {
	r.inserted = append(r.inserted, cc)
	return cc, nil
}

type stubMakeupClassRepository struct {
	existing map[string]struct{}
	inserted []academicdomain.MakeupClass
}

func (r *stubMakeupClassRepository) Exists(_ context.Context, subjectID string, date time.Time, period string) (bool, error) {
	_, ok := r.existing[cancelledKey(subjectID, date, period)]
	return ok, nil
}

func (r *stubMakeupClassRepository) Insert(_ context.Context, mc academicdomain.MakeupClass) (academicdomain.MakeupClass, error) {
	r.inserted = append(r.inserted, mc)
	return mc, nil
}

type stubRoomChangeRepository struct {
	existing map[string]struct{}
	inserted []academicdomain.RoomChange
}

func (r *stubRoomChangeRepository) Exists(_ context.Context, subjectID string, date time.Time, period, originalRoomID, newRoomID string) (bool, error) {
	_, ok := r.existing[cancelledKey(subjectID, date, period)+"|"+originalRoomID+"|"+newRoomID]
	return ok, nil
}

func (r *stubRoomChangeRepository) Insert(_ context.Context, rc academicdomain.RoomChange) (academicdomain.RoomChange, error) {
	r.inserted = append(r.inserted, rc)
	return rc, nil
}

func newScrapeServiceForTest(html string) (*ClassChangeScrapeService, *stubCancelledClassRepository, *stubMakeupClassRepository, *stubRoomChangeRepository) {
	cancelled := &stubCancelledClassRepository{existing: map[string]struct{}{}}
	makeup := &stubMakeupClassRepository{existing: map[string]struct{}{}}
	roomChange := &stubRoomChangeRepository{existing: map[string]struct{}{}}
	svc := NewClassChangeScrapeService(
		&stubScraperClient{html: html},
		&stubSubjectRepository{subjects: []domain.SubjectRef{
			{ID: "subject-1", SyllabusID: "101", Name: "情報数学"},
		}},
		&stubRoomRepository{rooms: []domain.RoomRef{
			{ID: "room-1", Name: "R791"},
			{ID: "room-2", Name: "R792"},
		}},
		cancelled,
		makeup,
		roomChange,
	)
	return svc, cancelled, makeup, roomChange
}

func TestClassChangeScrapeServiceRun(t *testing.T) {
	svc, cancelled, makeup, roomChange := newScrapeServiceForTest(scrapeTestHTML)

	summary, err := svc.Run(context.Background(), "user", "password", 2026)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if summary.CancelledInserted != 1 || summary.CancelledSkipped != 1 {
		t.Errorf("cancelled inserted=%d skipped=%d, want 1/1", summary.CancelledInserted, summary.CancelledSkipped)
	}
	if summary.MakeupInserted != 1 {
		t.Errorf("makeup inserted=%d, want 1", summary.MakeupInserted)
	}
	if summary.RoomChangeInserted != 1 {
		t.Errorf("room change inserted=%d, want 1", summary.RoomChangeInserted)
	}
	if len(summary.UnmatchedSubjectNames) != 1 || summary.UnmatchedSubjectNames[0] != "未知の科目" {
		t.Errorf("unmatched subjects = %v, want [未知の科目]", summary.UnmatchedSubjectNames)
	}

	if len(cancelled.inserted) != 1 {
		t.Fatalf("cancelled inserted rows = %d, want 1", len(cancelled.inserted))
	}
	cc := cancelled.inserted[0]
	if cc.Subject.ID != "subject-1" || cc.Date != "2026-05-12" || cc.Period != academicdomain.PeriodPeriod3 {
		t.Errorf("unexpected cancelled class: %+v", cc)
	}

	if len(makeup.inserted) != 1 {
		t.Fatalf("makeup inserted rows = %d, want 1", len(makeup.inserted))
	}
	if len(roomChange.inserted) != 1 {
		t.Fatalf("room change inserted rows = %d, want 1", len(roomChange.inserted))
	}
	rc := roomChange.inserted[0]
	if rc.OriginalRoom.ID != "room-1" || rc.NewRoom.ID != "room-2" {
		t.Errorf("unexpected room change: %+v", rc)
	}
}

func TestClassChangeScrapeServiceRunIsIdempotent(t *testing.T) {
	svc, cancelled, makeup, roomChange := newScrapeServiceForTest(scrapeTestHTML)
	date := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	cancelled.existing[cancelledKey("subject-1", date, "Period3")] = struct{}{}
	makeupDate := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	makeup.existing[cancelledKey("subject-1", makeupDate, "Period2")] = struct{}{}
	roomChangeDate := time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC)
	roomChange.existing[cancelledKey("subject-1", roomChangeDate, "Period4")+"|room-1|room-2"] = struct{}{}

	summary, err := svc.Run(context.Background(), "user", "password", 2026)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if summary.CancelledDuplicated != 1 || summary.CancelledInserted != 0 {
		t.Errorf("cancelled duplicated=%d inserted=%d, want 1/0", summary.CancelledDuplicated, summary.CancelledInserted)
	}
	if summary.MakeupDuplicated != 1 || summary.MakeupInserted != 0 {
		t.Errorf("makeup duplicated=%d inserted=%d, want 1/0", summary.MakeupDuplicated, summary.MakeupInserted)
	}
	if summary.RoomChangeDuplicated != 1 || summary.RoomChangeInserted != 0 {
		t.Errorf("room change duplicated=%d inserted=%d, want 1/0", summary.RoomChangeDuplicated, summary.RoomChangeInserted)
	}
	if len(cancelled.inserted)+len(makeup.inserted)+len(roomChange.inserted) != 0 {
		t.Errorf("expected no inserts on duplicated run")
	}
}
