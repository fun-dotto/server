// Package uuidv5 は RFC 9562 の名前ベース UUID (version 5, SHA-1) を生成する。
//
// 標準ライブラリの uuid パッケージは version 4 / version 7 の生成のみを提供し、
// 名前ベース UUID の生成関数を持たない。決定論的な ID を既存データと同じ値で
// 生成し続けるために、ここで RFC 9562 §5.5 に沿った実装を用意している。
package uuidv5

import (
	"crypto/sha1"
	"uuid"
)

// NamespaceURL は RFC 9562 が定義する URL 名前空間の UUID。
var NamespaceURL = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")

// NewSHA1 は namespace と name から version 5 の UUID を生成する。
// 同じ引数に対して常に同じ UUID を返す。
func NewSHA1(namespace uuid.UUID, name []byte) uuid.UUID {
	h := sha1.New()
	h.Write(namespace[:])
	h.Write(name)
	sum := h.Sum(nil)

	var u uuid.UUID
	copy(u[:], sum)
	u[6] = (u[6] & 0x0f) | 0x50 // version 5
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 9562 variant
	return u
}
