package uuidv5_test

import (
	"testing"
	"uuid"

	"github.com/fun-dotto/server/internal/shared/uuidv5"
)

// 期待値は移行前に github.com/google/uuid の NewSHA1 が生成していた値であり、
// 既存データの ID と一致し続けることを保証する。
func TestNewSHA1(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "cancelled class",
			key:  "urn:schedule-scripts:class-change:cancelled_class:abc123",
			want: "b2cd9c3e-9b7c-57b1-83d6-a55ec323f00d",
		},
		{
			name: "room change",
			key:  "urn:schedule-scripts:class-change:room_change:00000000-0000-0000-0000-000000000001",
			want: "dc103032-7555-568f-98cd-b3beadaa4bff",
		},
		{
			name: "makeup class",
			key:  "urn:schedule-scripts:class-change:makeup_class:xyz",
			want: "a04cc88b-5cc5-54da-bc49-a05fd0323671",
		},
		{
			name: "empty name",
			key:  "",
			want: "1b4db7eb-4057-5ddf-91e0-36dec72071f5",
		},
		{
			name: "RFC example URL",
			key:  "https://www.example.com/",
			want: "3d3ed9d2-aa3d-5fa6-90e8-ed662e90f559",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uuidv5.NewSHA1(uuidv5.NamespaceURL, []byte(tt.key)).String()
			if got != tt.want {
				t.Errorf("NewSHA1(%q) = %s, want %s", tt.key, got, tt.want)
			}
		})
	}
}

func TestNewSHA1IsDeterministic(t *testing.T) {
	name := []byte("urn:schedule-scripts:class-change:cancelled_class:abc123")
	if a, b := uuidv5.NewSHA1(uuidv5.NamespaceURL, name), uuidv5.NewSHA1(uuidv5.NamespaceURL, name); a != b {
		t.Errorf("同じ入力で異なる UUID が生成された: %v != %v", a, b)
	}
}

func TestNewSHA1SetsVersionAndVariant(t *testing.T) {
	u := uuidv5.NewSHA1(uuidv5.NamespaceURL, []byte("any"))
	if v := u[6] >> 4; v != 5 {
		t.Errorf("version = %d, want 5", v)
	}
	if variant := u[8] >> 6; variant != 0b10 {
		t.Errorf("variant bits = %b, want 10", variant)
	}
}

func TestNamespaceURL(t *testing.T) {
	want := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	if uuidv5.NamespaceURL != want {
		t.Errorf("NamespaceURL = %v, want %v", uuidv5.NamespaceURL, want)
	}
}
