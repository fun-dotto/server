package server

import "testing"

func TestAddr(t *testing.T) {
	tests := []struct {
		name string
		port string
		want string
	}{
		{
			name: "PORT 未設定なら :8080",
			port: "",
			want: ":8080",
		},
		{
			name: "PORT が設定されていればそのポートを使う (portless の割り当てポート)",
			port: "4123",
			want: ":4123",
		},
		{
			name: "PORT の前後の空白は無視する",
			port: " 4123 ",
			want: ":4123",
		},
		{
			name: "PORT が空白のみなら :8080",
			port: "   ",
			want: ":8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)

			if got := Addr(); got != tt.want {
				t.Errorf("Addr() = %q, want %q", got, tt.want)
			}
		})
	}
}
