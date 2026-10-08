---
root: false
targets: ["*"]
globs: ["**/*"]
---

# API デザインガイドライン

`proto/` 配下の Protobuf 定義は以下の規約に従う。`buf lint`（STANDARD）と `buf breaking`（FILE）を通過させること。

## 基本方針

- API 設計は原則として [Google API Improvement Proposals（AIP）](https://google.aip.dev/) に準拠する
- このガイドに書かれていない事項は AIP に従う
- AIP から外れる場合は、このガイドの「AIP からの逸脱」に理由とともに明記する。記載のない逸脱はしない
- 全体はリソース指向で設計する（[AIP-121](https://google.aip.dev/121)）

## ディレクトリ・パッケージ

- 配置: `proto/<domain>/<version>/<resource>.proto`（例: `proto/admin/v1/subject.proto`）
- `package` はディレクトリと一致させ、メジャーバージョンを末尾に付ける（例: `package admin.v1;`）（[AIP-185](https://google.aip.dev/185)、[AIP-191](https://google.aip.dev/191)）
- 1 ファイル 1 リソース。ファイル名は `lower_snake_case`（[AIP-191](https://google.aip.dev/191)）
- `go_package` は記述しない（`buf.gen.yaml` の managed mode で付与）

## ファイル構成

次の順で記述する。

1. `syntax = "proto3";`
2. `package`
3. `import`（アルファベット順）
4. 補助 message（リソースに従属する型）
5. リソース message
6. RPC ごとの Request / Response message（List → Get → Create → Update → Delete の順）
7. `service`

## 命名

[AIP-190](https://google.aip.dev/190)、[AIP-140](https://google.aip.dev/140) に従う。

- message / service / enum: `PascalCase`
- フィールド: `lower_snake_case`。前置詞を含めない
- `repeated` フィールドは複数形（[AIP-144](https://google.aip.dev/144)）
- 補助 message はリソース名を接頭辞にする（例: `SubjectFaculty`）
- service は `<Resource>Service`
- RPC は `<Verb><Resource>`（List のみ複数形: `ListSubjects`）（[AIP-131](https://google.aip.dev/131)〜[AIP-135](https://google.aip.dev/135)）
- Request / Response は RPC 名 + `Request` / `Response`。RPC 間で使い回さない
- enum の値は `UPPER_SNAKE_CASE` で enum 名を接頭辞にし、0 番は `<ENUM>_UNSPECIFIED` にする（[AIP-126](https://google.aip.dev/126)）

## 型

- 日時は `google.protobuf.Timestamp`、期間は `google.protobuf.Duration`（[AIP-142](https://google.aip.dev/142)）
- ID は `string`
- 未設定と既定値を区別する必要があるフィールドのみ `optional`（[AIP-149](https://google.aip.dev/149)）
- 取りうる値が決まっている文字列は enum にする（[AIP-126](https://google.aip.dev/126)）

## リソース message

- `id` をフィールド番号 1 に置く
- `created_at` / `updated_at` を末尾に置く（出力専用）

## 標準メソッド

| RPC      | AIP                                   | Request                                   | Response                                 |
| -------- | ------------------------------------- | ----------------------------------------- | ---------------------------------------- |
| `List`   | [AIP-132](https://google.aip.dev/132) | `page_size`, `page_token`                 | `repeated <Resource>`, `next_page_token` |
| `Get`    | [AIP-131](https://google.aip.dev/131) | `id`                                      | `<Resource>`                             |
| `Create` | [AIP-133](https://google.aip.dev/133) | `id` とタイムスタンプを除くフィールド     | `<Resource>`                             |
| `Update` | [AIP-134](https://google.aip.dev/134) | `id` + 更新対象フィールド + `update_mask` | `<Resource>`                             |
| `Delete` | [AIP-135](https://google.aip.dev/135) | `id`                                      | 空 message                               |

- `Update` の `update_mask` は `google.protobuf.FieldMask` で、フィールド番号 `99` に固定する（[AIP-134](https://google.aip.dev/134)、[AIP-161](https://google.aip.dev/161)）
  - `update_mask` が空の場合は、リクエストで値が入っているフィールドをすべて更新する
  - `update_mask` に存在しないフィールドが含まれる場合は `INVALID_ARGUMENT`
- 空の Response も省略せず定義する（例: `message DeleteSubjectResponse {}`）
- 標準メソッドで表せない操作はカスタムメソッドにする（[AIP-136](https://google.aip.dev/136)）

## ページネーション

カーソル方式を採用する（[AIP-158](https://google.aip.dev/158)）。

- List は最初からページネーションに対応させる。後から追加すると後方互換性を壊すため
- Request は `int32 page_size = 1;` と `string page_token = 2;` を持つ
- Response は `repeated <Resource>` を 1 番、`string next_page_token = 2;` を持つ
- `page_size` が 0 の場合はサーバー既定値（50）を使い、上限（1000）を超えた値は上限に丸める。負数は `INVALID_ARGUMENT`
- 返す件数が `page_size` より少なくても、それだけでは最終ページを意味しない。最終ページの判定は `next_page_token` が空文字かどうかで行う
- トークンは不透明な文字列とし、クライアントは中身を解釈・生成しない。サーバーはトークンの形式を予告なく変更してよい
- 2 ページ目以降は `page_size` 以外のリクエスト条件を最初と同じにする。条件が違う、または不正なトークンは `INVALID_ARGUMENT`
- ソートキーは一意になるようにする（末尾に `id` を加える）
- 総件数（`total_size`）は必要になるまで返さない

## エラー

[AIP-193](https://google.aip.dev/193) に従う。Connect のエラーとして返し、Response message にエラー用フィールドは持たせない。

### ステータスコード

[`google.rpc.Code`](https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto)（Connect では `connect.Code`）のうち最も具体的なものを使う。

| 状況                                                                       | Code                  |
| -------------------------------------------------------------------------- | --------------------- |
| リクエストの値が不正（必須項目の欠落、形式違反、不正な `page_token` など） | `INVALID_ARGUMENT`    |
| 対象リソースが存在しない                                                   | `NOT_FOUND`           |
| 作成しようとしたリソースが既に存在する                                     | `ALREADY_EXISTS`      |
| リソースの状態が操作の前提を満たさない                                     | `FAILED_PRECONDITION` |
| 認証情報がない、または無効                                                 | `UNAUTHENTICATED`     |
| 認証済みだが権限がない                                                     | `PERMISSION_DENIED`   |
| 同時更新の競合                                                             | `ABORTED`             |
| 一時的に処理できない（依存先の障害など。リトライ可能）                     | `UNAVAILABLE`         |
| 想定外のサーバー内部エラー                                                 | `INTERNAL`            |

- 権限のないリソースについては、存在しない場合でも `PERMISSION_DENIED` を返し、存在の有無を漏らさない（[AIP-211](https://google.aip.dev/211)）
- `UNKNOWN` は使わない

### メッセージ

- 開発者向けの英語で書き、原因と対処がわかるようにする
- 内部実装の詳細（SQL、スタックトレース、内部ホスト名など）を含めない
- エンドユーザー向けの文言が必要な場合は `google.rpc.LocalizedMessage` を details に付ける

### エラー詳細（details）

- すべてのエラーに `google.rpc.ErrorInfo` を 1 つ付ける
  - `reason`: エラーの種類を表す `UPPER_SNAKE_CASE` の識別子（最大 63 文字。例: `SUBJECT_NOT_FOUND`）。一度公開したら変更しない
  - `domain`: `<domain>.dotto`（`<domain>` は proto のドメイン。例: `admin.dotto`）
  - `metadata`: 原因の特定に必要な値（例: `{"subject_id": "..."}`）。キーは `lowerCamelCase`
- `INVALID_ARGUMENT` には `google.rpc.BadRequest` を付け、不正なフィールドごとに `field_violations` を入れる。`field` はリクエスト内のパスを `lower_snake_case` のドット区切りで書く（例: `faculties.0.faculty_id`）
- リトライ可能なエラーには必要に応じて `google.rpc.RetryInfo` を付ける
- クライアントの分岐は `Code` と `ErrorInfo.reason` で行い、メッセージ文字列に依存させない

```go
err := connect.NewError(connect.CodeNotFound, errors.New("subject not found"))
if d, derr := connect.NewErrorDetail(&errdetails.ErrorInfo{
	Reason:   "SUBJECT_NOT_FOUND",
	Domain:   "admin.dotto",
	Metadata: map[string]string{"subjectId": id},
}); derr == nil {
	err.AddDetail(d)
}
return nil, err
```

## 互換性

[AIP-180](https://google.aip.dev/180) に従い、同一メジャーバージョン内では後方互換性を保つ。

- 公開済みフィールドの番号・型・名前は変更しない
- フィールドを削除する場合は `reserved` で番号と名前を予約する
- 既存フィールドの意味やデフォルトの挙動を変えない
- 互換性を壊す変更は新しいメジャーバージョン（`v2`）で行う（[AIP-185](https://google.aip.dev/185)）

## AIP からの逸脱

既存定義との整合性、または `buf lint`（STANDARD）の規則を優先するため、以下は AIP と異なる方式を採る。ドメイン固有の逸脱は各ドメインのガイドライン（例: `admin-api-design.md`）に記載する。

| 項目                              | AIP                                                                                                                                                    | 本プロジェクト                              | 理由                                                                       |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------- | -------------------------------------------------------------------------- |
| リソースの識別子                  | `name` にリソース名（`subjects/{subject}`）（[AIP-122](https://google.aip.dev/122)）                                                                   | `id` に ID 文字列                           | 既存 REST API・DB との整合性                                               |
| タイムスタンプのフィールド名      | `create_time` / `update_time`（[AIP-148](https://google.aip.dev/148)）                                                                                 | `created_at` / `updated_at`                 | 既存 REST API・DB との整合性                                               |
| Get / Create / Update の Response | リソース message をそのまま返す（[AIP-131](https://google.aip.dev/131)、[AIP-133](https://google.aip.dev/133)、[AIP-134](https://google.aip.dev/134)） | `<Rpc>Response` でリソースを包む            | `buf lint` の `RPC_RESPONSE_STANDARD_NAME` / `RPC_REQUEST_RESPONSE_UNIQUE` |
| Delete の Response                | `google.protobuf.Empty`（[AIP-135](https://google.aip.dev/135)）                                                                                       | 空の `<Rpc>Response`                        | 同上。将来のフィールド追加にも対応できる                                   |
| Create / Update の Request        | リソース message をフィールドとして持つ（[AIP-133](https://google.aip.dev/133)、[AIP-134](https://google.aip.dev/134)）                                | リソースのフィールドを Request に直接並べる | 既存定義との整合性                                                         |

## コード生成

- 定義を変更したら `buf generate` を実行し、`gen/` の生成物を同じブランチでコミットする
- `gen/` は手で編集しない
