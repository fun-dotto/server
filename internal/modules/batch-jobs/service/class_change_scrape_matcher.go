package service

import (
	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
)

// subjectMatcher は科目名から subjects.id を引く。
// キーは NFKC 正規化した科目名と、さらに「（旧:…）」注釈を除去した科目名の 2 種類。
// 同一キーが衝突した場合は先勝ちする。
type subjectMatcher struct {
	byName map[string]string
}

func newSubjectMatcher(subjects []domain.SubjectRef) *subjectMatcher {
	byName := make(map[string]string, len(subjects)*2)
	register := func(key, id string) {
		if key == "" {
			return
		}
		if _, ok := byName[key]; !ok {
			byName[key] = id
		}
	}
	for _, s := range subjects {
		normalized := NormalizeSubjectName(s.Name)
		register(normalized, s.ID)
		register(NormalizeSubjectName(StripLegacyAnnotation(normalized)), s.ID)
	}
	return &subjectMatcher{byName: byName}
}

func (m *subjectMatcher) resolve(lessonName string) (string, bool) {
	normalized := NormalizeSubjectName(lessonName)
	if id, ok := m.byName[normalized]; ok {
		return id, true
	}
	if id, ok := m.byName[NormalizeSubjectName(StripLegacyAnnotation(normalized))]; ok {
		return id, true
	}
	return "", false
}

// roomMatcher は教室名から rooms.id を引く。キーは NormalizeRoomName で正規化した教室名。
type roomMatcher struct {
	byName map[string]string
}

func newRoomMatcher(rooms []domain.RoomRef) *roomMatcher {
	byName := make(map[string]string, len(rooms))
	for _, r := range rooms {
		key := NormalizeRoomName(r.Name)
		if key == "" {
			continue
		}
		if _, ok := byName[key]; !ok {
			byName[key] = r.ID
		}
	}
	return &roomMatcher{byName: byName}
}

func (m *roomMatcher) resolve(roomName string) (string, bool) {
	id, ok := m.byName[NormalizeRoomName(roomName)]
	return id, ok
}
