package event

import (
	"testing"

	"github.com/fun-dotto/server/internal/modules/academic/assets"
)

func TestLoadTermsAssets(t *testing.T) {
	terms, err := LoadTerms(assets.TermsJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 10 {
		t.Fatalf("got %d terms, want 10", len(terms))
	}
}
