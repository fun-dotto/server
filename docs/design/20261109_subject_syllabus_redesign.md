# 科目テーブル・シラバステーブル再設計

- 作成日: 2026-11-09
- ステータス: Draft

## 背景

現状の `syllabuses` は教務システムのシラバスID（text）を主キーとし、`subjects.syllabus_id` がそれを UNIQUE 付き FK で参照している（実質 1:1）。
`name` / `credit` / 教員名 / 学年などがシラバスと科目の双方に存在し、どちらが正本か曖昧になっている。

## 基本方針

- **シラバステーブル**: [教務システムのシラバス](https://students.fun.ac.jp/lesson/SyllabusList) の一覧・詳細ページから取得できる情報を **網羅的に、構造化せず文字列のまま** 保存する（取得元の写し）。
- **科目テーブル**: シラバステーブルを元に生成し、**構造化したデータ** を保存する（アプリが利用する正本）。
- シラバスの `id` は教務システムの ID から `lesson_id` に改名し、`id` は UUID で採番し直す。
- `subjects.syllabus_id` は新しい UUID の `syllabuses.id` を参照する。

## 変更点

### 1. DB スキーマ（Atlas マイグレーション追加）

#### `syllabuses`

- 既存 `id`（text）→ `lesson_id` に改名し `NOT NULL` + ユニーク制約。
- 新規 `id uuid PRIMARY KEY DEFAULT gen_random_uuid()`。
- ユニークキーは lesson_id が年度を跨いで同一かにより `lesson_id` 単独 or `(lesson_id, year)`（**要確認**）。
- 一覧ページ由来の列を追加（推測。HTML を確認して確定する）:
  - 科目名（現 `name` は授業名で別項目）
  - 年度
  - 開講期
  - 履修可能（開始学年）／履修可能（終了学年）
  - 所属
  - その他一覧にのみ表示される項目
- 全列を `text` で保持。現 `credit bigint` も方針に合わせ `text` 化を検討。
- 取得日時（`fetched_at`）列の追加を検討。

#### `subjects`

- `syllabus_id`: `text` → `uuid`、新 `syllabuses.id` を参照。
- `syllabus_id` の UNIQUE（1:1）を維持するか、1 シラバス → 複数科目を許容するか（**要確認**）。
- 構造化列の整理:
  - `name` / `credit` 等はシラバスと重複しても科目側に構造化済みの正本として保持。
  - 学年・クラス → `subject_eligible_attributes`
  - 必修/選択 → `subject_requirements`
  - 教員 → `subject_faculties`
- 現在 `syllabus_id` を表示順に使っているため、代替の並び順キー（科目コード、または `lesson_id` を `subjects` に持たせる等）を追加する。

#### データ移行 SQL

1. `syllabuses` に UUID を採番し、旧 `id` を `lesson_id` にコピー。
2. `subjects` に新列を追加し、旧 `syllabus_id`（= lesson_id）で JOIN して新 UUID を設定。
3. FK・ユニークインデックスを張り替え、旧列を削除。

### 2. server（Go）

- `internal/shared/model/syllabus.go` / `subject.go`: GORM タグ・型変更。
- `internal/modules/academic/{domain,repository}`（`syllabus.go`, `convert.go`, `subject.go`）: `GetByID` の型変更、lesson_id 検索の追加。
- `Order("...syllabus_id ASC")` の置換（UUID 順は無意味なため新しい並び順キーへ）:
  - `internal/modules/academic/repository/`: `room_change.go`, `timetable_items.go`, `subject.go`, `cancelled_class.go`, `makeup_class.go`
  - `internal/modules/appapi/handler/`: `timetable_item.go`, `subject.go`, `course_registration.go`, `class_change.go`
- `internal/modules/adminapi/handler/syllabus.go`: キーを UUID に変更し、lesson_id も扱う。
- `internal/modules/adminapi/handler/subject.go`: `syllabus_id` の型・バリデーション変更。
- **新規ジョブ**（現リポジトリにスクレイパー・変換処理は存在しない）:
  - `cmd/*-job`: 教務システムから取得し `syllabuses` に upsert するジョブ。
  - `cmd/*-job`: `syllabuses` をパースして `subjects` と子テーブルを生成・更新するジョブ。手動修正した科目データを上書きしないためのルールが必要。

### 3. API 定義

- `proto/admin/v1/syllabus.proto`: `id` を UUID に、`lesson_id` と追加列を追加。
- `proto/admin/v1/subject.proto`: `syllabus_id` の説明更新。
- `proto/app/v1/subject.proto` `SubjectSyllabus.id`: 現コメントは「教務システムのシラバスID」。UUID 化するなら `lesson_id` フィールドを追加。
- `api/openapi/academic/openapi.yaml`: `Syllabus` スキーマ、`/v1/subjects/{id}/syllabus`。
- `gen/` の再生成。

### 4. app（Flutter）

- **互換性リスク最大**: `subject_detail_screen.dart` / `subject_review_new_screen.dart` で `syllabus.id` をレビュー（フィードバック）の `lessonId` として使用している。`id` が UUID になると既存レビューとの紐付けが切れるため、
  1. API に `lesson_id` を追加
  2. app の参照を `syllabus.lessonId` に切り替え
  3. `id` の意味を変更

  の順で段階移行する。
- `apps/dotto/assets/syllabus.db`（同梱 SQLite）とテスト用 fake データ（`test/helpers/fake_subject_data_source.dart` 等）の追従。
- リリース済みの旧バージョンのアプリが新 API を叩いた場合の互換性を考慮する。

## 未決事項

1. lesson_id は年度を跨いで同一か（`syllabuses` のユニークキー）。
2. 1 シラバス → 複数科目を許容するか（`subjects.syllabus_id` の UNIQUE）。
3. 一覧・詳細ページの全項目（HTML サンプルを入手して列を確定）。
4. 今回のスコープ（マイグレーション + server のみ／取得・生成ジョブ／app 対応まで）。
