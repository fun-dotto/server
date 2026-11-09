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

#### 取得元ページの調査結果（2026-11-09）

ログイン不要で取得できる。

- 一覧: `https://students.fun.ac.jp/lesson/SyllabusList`（ASP.NET WebForms。`__VIEWSTATE` 等を付けて POST で検索・ページング。2026 年度は 341 件、1 ページ 20 件）
- 詳細: `https://students.fun.ac.jp/lesson/Syllabus?lesson_id={授業コード}&year={年度}`
- **同じ `lesson_id` が年度をまたいで使われる**（例: `100101` 認知科学は 2025・2026 の両方に存在し、担当教員が異なる）。

一覧の列:

| 表示名 | 備考 |
| --- | --- |
| 履修年度 | |
| 授業コード | = `lesson_id` |
| 授業科目名 | 詳細の「授業名」と同じ値 |
| Course | 詳細の「Course」と同じ値 |
| 担当教員名 | 主担当のみ（詳細は全員 + `◎`） |
| 授業形態 | |
| 表示順 | 一覧にのみ存在 |

詳細ページの項目（`*` は科目によって表示されない項目）:

- 基本情報: 年度、授業名、Course、配当年次、開講期、単位数、担当教員名（改行区切り、主担当に `◎`）、実務家教員区分、複数人担当形式、授業形態
- 本文: 授業の概要、授業の到達目標、提出課題等、成績の評価方法・基準、テキスト、参考書、履修条件、事前学習、事後学習、履修上の留意点、キーワード/keywords、対象コース・領域、科目群・科目区分、DSOP対象科目`*`、教授言語`*`、授業・試験の形式`*`、授業内容とスケジュール
- 英語欄`*`（主に大学院）: (E)Course outline、(E)Course Objectives、(E)Course Schedule、(E)Assessment、(E)Language of Instruction、(E)Requirements for registration
- ディプロマポリシー: ディプロマポリシー（学部／大学院の DP 一覧）、ディプロマポリシーとの対応（例: `DP2,DP4`）、ディプロマポリシーの各領域と配点（DP1〜DP10 の各枠に配点と S/A/B/C/D の評価基準）

#### `syllabuses`

取得元の写しとして、値は加工せず文字列のまま保存する（空欄・非表示は `''`）。

- 主キー: `id uuid DEFAULT gen_random_uuid()`
- **バージョン管理**: 取得した内容が前回と異なる場合は既存行を更新せず、新しい行（新しい `id`）として追加する。同じ `(lesson_id, year)` に複数バージョンが存在しうる。
- どのバージョンを使うかは `subjects.syllabus_id` の紐付けで決まる。
- 内容比較のため `content_hash`（取得値から計算したハッシュ）を持ち、`(lesson_id, year, content_hash)` をユニークキーにする。
- インデックス: `(lesson_id, year, fetched_at)`（最新バージョンの検索用）

| 列 | 型 | 取得元 | 既存列からの変更 |
| --- | --- | --- | --- |
| `id` | uuid | — | 新規（UUID 採番） |
| `lesson_id` | text | 授業コード | 旧 `id` を改名 |
| `year` | bigint | 年度 | 新規（キーのため数値） |
| `display_order` | text | 一覧: 表示順 | 新規 |
| `list_staff_name` | text | 一覧: 担当教員名 | 新規 |
| `name` | text | 授業名 | |
| `en_name` | text | Course | |
| `grades` | text | 配当年次 | |
| `term` | text | 開講期 | 新規 |
| `credit` | text | 単位数（例: `2単位`） | bigint → text |
| `faculty_names` | text | 担当教員名 | |
| `practical_home_faculty_category` | text | 実務家教員区分 | |
| `multiple_person_teaching_form` | text | 複数人担当形式 | |
| `teaching_form` | text | 授業形態 | |
| `summary` | text | 授業の概要 | |
| `learning_outcomes` | text | 授業の到達目標 | |
| `assignments` | text | 提出課題等 | |
| `evaluation_method` | text | 成績の評価方法・基準 | |
| `textbooks` | text | テキスト | |
| `reference_books` | text | 参考書 | |
| `prerequisites` | text | 履修条件 | |
| `pre_learning` | text | 事前学習 | |
| `post_learning` | text | 事後学習 | |
| `notes` | text | 履修上の留意点 | |
| `keywords` | text | キーワード/keywords | |
| `target_courses` | text | 対象コース・領域 | |
| `classifications` | text | 科目群・科目区分 | |
| `dsop_subject` | text | DSOP対象科目 | |
| `teaching_language` | text | 教授言語 | |
| `teaching_and_exam_form` | text | 授業・試験の形式 | |
| `contents_and_schedule` | text | 授業内容とスケジュール | |
| `en_summary` | text | (E)Course outline | 新規 |
| `en_learning_outcomes` | text | (E)Course Objectives | 新規 |
| `en_contents_and_schedule` | text | (E)Course Schedule | 新規 |
| `en_evaluation_method` | text | (E)Assessment | 新規 |
| `en_teaching_language` | text | (E)Language of Instruction | 新規 |
| `en_prerequisites` | text | (E)Requirements for registration | 新規 |
| `diploma_policies` | text | ディプロマポリシー | 新規 |
| `diploma_policy_correspondence` | text | ディプロマポリシーとの対応 | 新規 |
| `diploma_policy_rubrics` | jsonb | 各領域と配点（DP1〜10 × 配点・S〜D の文字列） | 新規 |
| `content_hash` | text | 取得値から計算 | 新規 |
| `fetched_at` | timestamptz | 取得日時（同一内容を再取得した場合は最終確認日時として更新） | 新規 |
| `created_at` / `updated_at` | timestamptz | — | |

`diploma_policy_rubrics` は 60 個の値になるため列には展開せず、`[{"weight":"0","s":"...","a":"...","b":"...","c":"...","d":"..."}, ...]` の形で文字列のまま保存する。

#### `subjects`

- `syllabus_id`: `text` → `uuid`、新 `syllabuses.id` を参照。**UNIQUE は維持する**（1 つのシラバスバージョンは 1 科目からのみ参照される）。科目と授業の対応は、紐付いたバージョンの `(lesson_id, year)` で決まる。
- 構造化した正本として以下を持つ（シラバスと重複しても保持する）。

| 列・子テーブル | 生成元（`syllabuses`） | 変換 |
| --- | --- | --- |
| `name` | `name` | そのまま |
| `year` | `year` | そのまま |
| `semester` | `term` | `前期`/`後期`/`1Q`〜`4Q`/`通年` → enum |
| `credit` | `credit` | `2単位` → `2` |
| `classification` | `classifications` | 必修/選択などを enum に |
| `cultural_subject_category` | `classifications` | 教養科目群の区分を enum に |
| `subject_faculties` | `faculty_names` | 改行で分割し、`◎` を `is_primary`。`faculties` と名前で突合 |
| `subject_eligible_attributes` | `grades`、`name` の末尾（例: `1-EFGH`, `2-GHI`, `1～4`） | 学年・クラスに分解 |
| `subject_requirements` | `target_courses`、`classifications` | コース × 必修/選択 |

- 異なる `(lesson_id, year)` のバージョンへの付け替えは許可しない（admin-api でバリデーション）。
- 表示順: 現在 `syllabus_id`（= 授業コード）で並べている箇所は、`syllabuses` を JOIN して `syllabuses.lesson_id` で並べる。

#### データ移行 SQL

1. `syllabuses` に `id uuid` を追加して採番し、旧 `id` を `lesson_id` に改名する。
2. 既存行の `year` は、紐づく `subjects.year` から埋める。
3. `subjects` に新しい `syllabus_id uuid` を追加し、旧 `syllabus_id`（= lesson_id）で JOIN して設定する。
4. FK・ユニークインデックスを張り替え、旧列を削除する。
5. `credit` を text に変換する（`2` → `2単位`）。

### 2. server（Go）

- `internal/shared/model/syllabus.go` / `subject.go`: GORM タグ・型変更。
- `internal/modules/academic/{domain,repository}`（`syllabus.go`, `convert.go`, `subject.go`）: `GetByID` の型変更。
- `Order("...syllabus_id ASC")` を `syllabuses.lesson_id` での並びに置換:
  - `internal/modules/academic/repository/`: `room_change.go`, `timetable_items.go`, `subject.go`, `cancelled_class.go`, `makeup_class.go`
  - `internal/modules/appapi/handler/`: `timetable_item.go`, `subject.go`, `course_registration.go`, `class_change.go`
- `internal/modules/adminapi/handler/syllabus.go`: キーを UUID に変更し、`lesson_id`・`year` と追加列を扱う。
- `internal/modules/adminapi/handler/subject.go`: `syllabus_id` の型・バリデーション変更。
- 新規ジョブ:
  - **シラバス取得ジョブ**（`cmd/fetch-syllabuses-job`）: 年度を指定して一覧を全ページ取得し、各詳細ページを取得する。`content_hash` が既存バージョンと一致すれば `fetched_at` のみ更新し、異なれば新しいバージョンとして INSERT する。アクセス間隔を空ける。
    - 新しいバージョンを追加した授業のうち、既存の科目が古いバージョンに紐付いているものを集計し、**差分（授業コード・年度・授業名・変更された項目名）を Slack に通知する**。付け替えは通知を見て admin-api から手動で行う。
    - Slack への送信は Incoming Webhook を使い、URL は環境変数（Secret Manager 経由）で渡す。現在リポジトリに Slack 連携はないため新規に追加する。
  - **科目生成ジョブ**（`cmd/build-subjects-job`）: 年度を指定して実行する。
    - 科目が存在しない `(lesson_id, year)` は、最新バージョンから `subjects` と子テーブルを新規作成する。
    - 既存の科目は更新しない（紐付いているバージョンを使い続ける）。新しいバージョンへの付け替えは admin-api で行い、付け替え時に紐付け先のシラバスから構造化データを再生成する。
    - 担当教員名が `faculties` に存在しない場合は、その教員を `subject_faculties` に入れずに警告ログを出す。
  - どちらも既存ジョブ（`build-class-change-notifications-job` など）と同じ形で `cmd/` に作り、`Dockerfile`・`mise.toml`・`.github/actions/set-env/action.yml` に追加して、同じ方法でデプロイ・スケジュール実行する。
  - スケジュール: **毎週日曜日**に、シラバス取得ジョブ → 科目生成ジョブの順で実行する（対象は当年度）。

### 3. API 定義

API ごとに互換性の要件が異なる。

| API | 利用者 | 方針 |
| --- | --- | --- |
| app-bff-api（`api/openapi/app`） | リリース済みアプリ | **スキーマ変更なし** |
| app-api（`proto/app/v1`） | 開発中のアプリ（未リリース） | 後方互換性を考慮せず変更してよい |
| academic-api（`api/openapi/academic`） | app-bff-api（内部） | app-bff-api が旧値を返せるよう変更 |
| admin-api（`proto/admin/v1`） | 管理画面 | 変更する |

#### app-bff-api（スキーマ変更なし）

- `AcademicService.Syllabus.id` は、リリース済みアプリがレビュー（フィードバック）のキー `lessonId` として使っている。**引き続き教務システムの ID（= `lesson_id`）を返す。**
- `internal/modules/app/external/subject_model.go` の `ToDomainSyllabus` で、`ID` に academic-api の `lessonId` を詰める。
- そのほか、レスポンスに含まれる列が型変更（例: `credit` の text 化）の影響を受ける場合は、BFF 側で旧型へ変換する。

#### academic-api

- `api/openapi/academic/openapi.yaml` の `Syllabus` スキーマ: `id` を UUID にし、`lessonId` と追加列を加える。
- `/v1/subjects/{id}/syllabus` はそのまま（subject の id は UUID のまま変わらない）。

#### app-api

- `proto/app/v1/subject.proto` の `SubjectSyllabus`: `id` を UUID にし、`lesson_id` と追加列を加える。互換性は考慮しない。

#### admin-api

- `proto/admin/v1/syllabus.proto`: `id` を UUID に、`lesson_id` と追加列を追加。
- `proto/admin/v1/subject.proto`: `syllabus_id` の説明更新。

#### 共通

- `gen/` の再生成。

### 4. app（Flutter）

- 現在の実装は app-api を使っており、未リリースのため互換性は考慮しない。
- `subject_detail_screen.dart` / `subject_review_new_screen.dart` で `syllabus.id` をレビュー（フィードバック）の `lessonId` に使っている。`syllabus.lessonId` を参照するよう変更する（UUID を使うと既存レビューと紐付かなくなるため）。
- `apps/dotto/assets/syllabus.db`（同梱 SQLite）とテスト用 fake データ（`test/helpers/fake_subject_data_source.dart` 等）の追従。

## 決定事項

- `lesson_id` は年度をまたいで同じ値 → 授業の識別は `(lesson_id, year)`。
- `syllabuses` は複数バージョンを保持し、どのバージョンを使うかは `subjects.syllabus_id` の紐付けで決める。
- 担当教員名が `faculties` に存在しない場合はスキップして警告ログ。
- ジョブは既存ジョブと同じ形で作り、毎週日曜日に実行する。
- シラバスの差分は Slack に通知し、科目の付け替えは admin-api から手動で行う。
- シラバスと科目は 1:1 のまま。
- 今回の実装範囲: DB + server、シラバス取得ジョブ、科目生成ジョブ（app 対応は対象外）。

## 未決事項

なし。
