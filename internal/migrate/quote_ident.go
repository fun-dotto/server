package migrate

// quoteIdent は PostgreSQL の識別子をダブルクォートで囲み、
// 含まれるダブルクォートを二重化してエスケープする。
func quoteIdent(name string) string {
	quoted := make([]byte, 0, len(name)+2)
	quoted = append(quoted, '"')
	for i := 0; i < len(name); i++ {
		if name[i] == '"' {
			quoted = append(quoted, '"', '"')
		} else {
			quoted = append(quoted, name[i])
		}
	}
	quoted = append(quoted, '"')
	return string(quoted)
}
