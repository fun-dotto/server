package service

import "testing"

func TestNormalizeSubjectName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "全角英数はNFKCで半角化される", in: "ＡＩ基礎", want: "AI基礎"},
		{name: "連続する空白は1つに圧縮される", in: "情報　工学   概論", want: "情報 工学 概論"},
		{name: "前後の空白は除去される", in: "  データ構造  ", want: "データ構造"},
		{name: "大文字小文字は区別する", in: "Python入門", want: "Python入門"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeSubjectName(tt.in); got != tt.want {
				t.Errorf("NormalizeSubjectName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestStripLegacyAnnotation(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "全角括弧の旧科目名注釈を除去する", in: "情報科学演習（旧:情報処理演習）", want: "情報科学演習"},
		{name: "全角コロンの注釈を除去する", in: "情報科学演習（旧：情報処理演習）", want: "情報科学演習"},
		{name: "半角括弧の注釈を除去する", in: "情報科学演習(旧:情報処理演習)", want: "情報科学演習"},
		{name: "注釈が無ければそのまま", in: "情報科学演習", want: "情報科学演習"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripLegacyAnnotation(tt.in); got != tt.want {
				t.Errorf("StripLegacyAnnotation(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeRoomName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "全角英数字はNFKCで半角化される", in: "Ｒ７８１", want: "r781"},
		{name: "全角スペースは通常の空白として圧縮される", in: "大　講義室", want: "大 講義室"},
		{name: "大文字小文字を無視する", in: "R781", want: "r781"},
		{name: "前後の空白は除去される", in: "  講堂  ", want: "講堂"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeRoomName(tt.in); got != tt.want {
				t.Errorf("NormalizeRoomName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
