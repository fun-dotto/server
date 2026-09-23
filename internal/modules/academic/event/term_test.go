package event

import (
	"os"
	"testing"
)

func TestLoadTermsAssets(t *testing.T) {
	b, err := os.ReadFile("../assets/terms.json")
	if err != nil {
		t.Fatal(err)
	}
	terms, err := LoadTerms(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 10 {
		t.Fatalf("got %d terms, want 10", len(terms))
	}
}
