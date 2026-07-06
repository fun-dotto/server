package domain

// SubjectRef は科目名照合に使う subjects テーブルの最小表現。
type SubjectRef struct {
	ID         string
	SyllabusID string
	Name       string
}
