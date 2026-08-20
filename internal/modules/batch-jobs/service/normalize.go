package service

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// legacyAnnotationPattern はポータル側の科目名に付く「旧科目名」注釈を検出する。
var legacyAnnotationPattern = regexp.MustCompile(`（旧[:：][^）]*）|\(旧[:：][^)]*\)`)

// NormalizeSubjectName は科目名照合用のキー: NFKC 正規化 + 連続空白の圧縮 + 前後空白除去。
// 大文字小文字は区別する（Python 版 normalize_for_match と同等）。
func NormalizeSubjectName(s string) string {
	t := norm.NFKC.String(s)
	return strings.TrimSpace(strings.Join(strings.Fields(t), " "))
}

// StripLegacyAnnotation は科目名に付く「（旧:…）」注釈をすべて除去し前後空白を整える。
func StripLegacyAnnotation(s string) string {
	return strings.TrimSpace(legacyAnnotationPattern.ReplaceAllString(s, ""))
}

// NormalizeRoomName は教室名照合用のキー: NFKC 正規化 + 連続空白の圧縮 + 前後空白除去 + 大文字小文字無視。
func NormalizeRoomName(s string) string {
	t := norm.NFKC.String(s)
	t = strings.Join(strings.Fields(t), " ")
	return strings.ToLower(strings.TrimSpace(t))
}
