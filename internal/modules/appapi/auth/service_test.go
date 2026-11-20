package auth

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/api/idtoken"
)

func TestGoogleServiceVerifier(t *testing.T) {
	v, err := NewGoogleServiceVerifier("https://app-api.example.com", " web@example.iam.gserviceaccount.com , ")
	if err != nil {
		t.Fatal(err)
	}
	v.validate = func(_ context.Context, token, audience string) (*idtoken.Payload, error) {
		if audience != "https://app-api.example.com" {
			return nil, errors.New("audience mismatch")
		}
		switch token {
		case "web":
			return &idtoken.Payload{Claims: map[string]any{"email": "web@example.iam.gserviceaccount.com", "email_verified": true}}, nil
		case "other":
			return &idtoken.Payload{Claims: map[string]any{"email": "other@example.iam.gserviceaccount.com", "email_verified": true}}, nil
		case "unverified":
			return &idtoken.Payload{Claims: map[string]any{"email": "web@example.iam.gserviceaccount.com"}}, nil
		}
		return nil, errors.New("invalid")
	}

	tests := []struct {
		token   string
		wantErr bool
	}{
		{"web", false},
		{"other", true},
		{"unverified", true},
		{"invalid", true},
	}
	for _, tt := range tests {
		if err := v.Verify(context.Background(), tt.token); (err != nil) != tt.wantErr {
			t.Errorf("Verify(%q) err = %v, wantErr %v", tt.token, err, tt.wantErr)
		}
	}

	if _, err := NewGoogleServiceVerifier("", "a@example.com"); err == nil {
		t.Error("empty audience should be rejected")
	}
	if _, err := NewGoogleServiceVerifier("https://x", " , "); err == nil {
		t.Error("empty accounts should be rejected")
	}
}
