package model

import "testing"

func TestSyllabusContentHash(t *testing.T) {
	a := Syllabus{LessonID: "1", Year: 2026, Name: "A", Summary: "x"}
	b := a
	if a.ComputeContentHash() != b.ComputeContentHash() {
		t.Fatal("same content must have the same hash")
	}
	b.Summary = "y"
	if a.ComputeContentHash() == b.ComputeContentHash() {
		t.Fatal("different content must have different hashes")
	}
	changed := ChangedSyllabusFields(&a, &b)
	if len(changed) != 1 || changed[0].Column != "summary" {
		t.Fatalf("changed = %+v", changed)
	}
}
