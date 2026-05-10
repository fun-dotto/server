// Package openapispec は academic API の OpenAPI 仕様を distroless で読めるよう
// バイナリに同梱するためだけのパッケージ。Cloud Run Service が起動するときに
// kin-openapi の LoadFromData に渡してリクエストバリデータを構築する。
package openapispec

import _ "embed"

// Spec は api/openapi/academic/openapi.yaml をバイナリへ埋め込んだバイト列。
//
//go:embed openapi.yaml
var Spec []byte
