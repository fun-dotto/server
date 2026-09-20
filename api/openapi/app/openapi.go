// Package openapispec は app API の OpenAPI 仕様を distroless で読めるよう
// バイナリに同梱するためだけのパッケージ。仕様の原本は同ディレクトリの
// openapi.yaml だけに置き、コピーは持たない。
package openapispec

import _ "embed"

// Spec は openapi.yaml をバイナリへ埋め込んだバイト列。
//
//go:embed openapi.yaml
var Spec []byte
