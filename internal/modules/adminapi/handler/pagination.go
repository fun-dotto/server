package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

const (
	defaultPageSize = 50
	maxPageSize     = 1000
)

func normalizePageSize(size int32) (int, error) {
	switch {
	case size < 0:
		return 0, invalidArgument("INVALID_PAGE_SIZE", "page_size", "must not be negative")
	case size == 0:
		return defaultPageSize, nil
	case size > maxPageSize:
		return maxPageSize, nil
	default:
		return int(size), nil
	}
}

// encodePageToken はソートキーの値を不透明なトークンにする。
func encodePageToken(cursor []any) (string, error) {
	b, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func decodePageToken(token string, wantLen int) ([]any, error) {
	if token == "" {
		return nil, nil
	}
	invalid := invalidArgument("INVALID_PAGE_TOKEN", "page_token", "must be a next_page_token returned by a previous List call")
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, invalid
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var cursor []any
	if err := dec.Decode(&cursor); err != nil || len(cursor) != wantLen {
		return nil, invalid
	}
	for i, v := range cursor {
		switch v := v.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return nil, invalid
			}
			cursor[i] = n
		case string:
		default:
			return nil, invalid
		}
	}
	return cursor, nil
}
