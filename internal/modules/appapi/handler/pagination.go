package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

const (
	defaultPageSize = 50
	maxPageSize     = 1000
)

// listRequest は page_size / page_token を持つ List の Request。
type listRequest interface {
	proto.Message
	GetPageSize() int32
	GetPageToken() string
}

type pageToken struct {
	Offset int    `json:"o"`
	Filter string `json:"f"`
}

// paginate は取得済みの一覧から 1 ページ分を切り出す。
// 上流のサービスがページングに対応していないため、オフセットをトークンに含める。
// トークンには絞り込み条件のハッシュを含め、条件が変わった場合は INVALID_ARGUMENT とする。
func paginate[T any](req listRequest, items []T) ([]T, string, error) {
	size, err := normalizePageSize(req.GetPageSize())
	if err != nil {
		return nil, "", err
	}
	filter, err := filterHash(req)
	if err != nil {
		return nil, "", err
	}
	offset := 0
	if t := req.GetPageToken(); t != "" {
		token, ok := decodePageToken(t)
		if !ok || token.Filter != filter || token.Offset < 0 {
			return nil, "", invalidArgument("INVALID_PAGE_TOKEN", "page_token", "must be a next_page_token returned by a previous List call with the same conditions")
		}
		offset = token.Offset
	}
	if offset >= len(items) {
		return []T{}, "", nil
	}
	end := min(offset+size, len(items))
	next := ""
	if end < len(items) {
		next = encodePageToken(pageToken{Offset: end, Filter: filter})
	}
	return items[offset:end], next, nil
}

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

// filterHash は page_size / page_token を除いた Request のハッシュを返す。
func filterHash(req listRequest) (string, error) {
	m := proto.Clone(req)
	fields := m.ProtoReflect().Descriptor().Fields()
	for _, name := range []string{"page_size", "page_token"} {
		if fd := fields.ByName(protoName(name)); fd != nil {
			m.ProtoReflect().Clear(fd)
		}
	}
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8]), nil
}

func encodePageToken(t pageToken) string {
	b, _ := json.Marshal(t)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodePageToken(s string) (pageToken, bool) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return pageToken{}, false
	}
	var t pageToken
	if err := json.Unmarshal(b, &t); err != nil {
		return pageToken{}, false
	}
	return t, true
}
