---
root: false
targets: ["*"]
globs: ["proto/admin/**/*", "internal/modules/admin/**/*"]
---

# Admin API デザインガイドライン

`proto/admin/` 配下の Admin API に固有のルール。共通ルールは `api-design.md` に従い、本ガイドはそれに追加・上書きする。

## エラー

- `google.rpc.ErrorInfo` の `domain` は `admin.dotto` とする

## AIP からの逸脱

`api-design.md` の「AIP からの逸脱」に加え、Admin API では以下を採る。

| 項目                               | AIP                                                                                                                               | Admin API                                                                                                                                                                                                                                                                                                                              | 理由                                                             |
| ---------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| リソースの識別子（GTFS・複合キー） | `id` をフィールド番号 1 に置く                                                                                                    | GTFS 由来のリソース（`Zone`、`Stop` など）は GTFS の列名（`zone_id`、`stop_id` など）を識別子とし、`FcmToken` は `token` を識別子とする。複合キーのリソース（`FacultyRoom`、`CourseRegistration`、`NotificationTargetUser`、`StopTime`、`CalendarDate`）はキーを構成するフィールドを並べ、Get / Delete の Request にもそれらを持たせる | GTFS の仕様・DB の主キーとの整合性                               |
| Create の Request の識別子         | サーバーが ID を採番し、クライアント指定の ID は `<resource>_id` で受け取る（[AIP-133](https://google.aip.dev/133)）              | ID が外部由来のリソース（`Syllabus`、`User`、GTFS 由来のリソースなど）は、Create の Request に識別子のフィールドを含める                                                                                                                                                                                                               | 外部システム（シラバス・認証基盤・GTFS）の ID をそのまま使うため |
| タイムスタンプの有無               | —                                                                                                                                 | DB に `created_at` / `updated_at` 列がないリソース（GTFS 由来のリソース、`User`、`NotificationTargetUser`）は省略する                                                                                                                                                                                                                  | DB との整合性                                                    |
| GTFS の時刻・日付                  | 時刻は `google.protobuf.Timestamp`（[AIP-142](https://google.aip.dev/142)）                                                       | GTFS 由来のリソース（`StopTime` など）は GTFS の形式の `string` で持つ                                                                                                                                                                                                                                                                 | GTFS の時刻は 24 時を超えうるなど `Timestamp` で表せないため     |
| 公開済みの名前                     | 前置詞・予約語・略語を避ける（[AIP-140](https://google.aip.dev/140)）、期間は `Duration`（[AIP-142](https://google.aip.dev/142)） | 公開済みの `available_from`、`notify_after`、`notify_before`、`class`、`COURSE_INFORMATION_*`、`android_ttl_seconds` は変更しない                                                                                                                                                                                                      | 後方互換性（[AIP-180](https://google.aip.dev/180)）              |
