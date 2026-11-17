# 過去問機能の API 化と GCS への移行

## 背景

過去問機能はサーバーを介さず、アプリ単体で動いている。

### 現状の構成（app リポジトリで確認、2026-11-17）

| 役割 | 実装 |
| --- | --- |
| 科目 → 過去問の紐付け | アプリにバンドルした `syllabus.db`（SQLite）の `detail` テーブル `過去問` 列（整数）。`SubjectDataSource.readPastExamId(lessonId)` で `LessonId` から引く |
| ファイル一覧 | Cloudflare R2 を S3 互換 API（`minio` パッケージ）で `listObjectsV2(prefix: pastExamId, recursive: true)`。`S3DataSource.getListObjectsKey` |
| ファイル取得 | 同じく `getObject(bucket, key)` でバイト列を取得し、一時ファイルに書き出して PDF ビューアで表示（`PdfDataSource.fetchPastExam`） |
| 表示名 | オブジェクトキーの最初の `/` 以降（`pastExamFileName`）。メタデータは持たない |
| 認可 | アプリ側で `isAuthenticated`（Google ログイン）を見て画面を出し分けているだけ |
| 資格情報 | `CLOUDFLARE_R2_ENDPOINT` / `ACCESS_KEY_ID` / `SECRET_ACCESS_KEY` / `BUCKET_NAME` を `--dart-define` でアプリに埋め込み |
| ルーティング | `/subjects/:id/past-exams/:pastExamId`、`/course/subjects/:id/past-exams/:pastExamId`。`pastExamId` に **オブジェクトキーそのもの** が入る |

### 課題

1. **R2 のアクセスキーがアプリバイナリに含まれている。** 抽出すればログインなしでバケット全体を列挙・取得できる（キーの権限次第では書き込みも）。ログイン必須はクライアント側の表示制御にすぎない。
2. **紐付けがアプリのリリースに縛られている。** 過去問の追加・修正のたびに R2 へのアップロードに加えて `syllabus.db` の更新とアプリリリースが必要。新しい科目は次のリリースまで過去問が出ない。
3. **`syllabus.db` の `LessonId` 依存。** サーバー側では科目が UUID の `subjects` に移行済み（[科目テーブル・シラバステーブル再設計](./20261109_subject_syllabus_redesign.md)）で、アプリの科目も API から取得するようになっていく。過去問だけが SQLite に残る。
4. **メタデータがない。** 年度・試験種別（中間 / 期末 / 小テスト等）・解答の有無などはファイル名頼みで、並び順や表示を制御できない。
5. **インフラが分散している。** 他のデータは GCP（Cloud SQL / Cloud Run）にあるのに、過去問だけ Cloudflare R2。権限管理・監査ログ・Terraform 管理が別系統になる。

## 方針

- **過去問のメタデータを Cloud SQL（`past_exams` テーブル）に持ち、科目（`subjects`）に紐付ける。** `syllabus.db` の `過去問` 列は使わない。
- **ファイルは GCS の非公開バケットに置く。** アプリは GCS に直接の資格情報を持たない。
- **App API（Connect RPC）に一覧取得とダウンロード URL 発行の RPC を追加する。** ダウンロードは **V4 署名付き URL（短期）** を返し、アプリはそれを HTTP GET する。ファイル本体を API サーバー経由で中継しない。
- **認可はサーバーで行う。** 過去問 RPC はログイン必須（`userRequiredProcedures`）にし、さらに大学ドメイン（`@fun.ac.jp`）のユーザーに限る。
- **登録・更新は Admin API から行う。** アップロードも署名付き URL（PUT）で GCS に直接送る。
- 移行期間はアプリの旧実装（R2）を残し、API 版のリリース後に R2 の資格情報を失効させる。

## 設計

### 1. GCS バケット（infra）

| 項目 | 値 |
| --- | --- |
| バケット | `swift2023groupc-past-exams-{prd,stg,dev}`（環境ごとに分離） |
| ロケーション | `asia-northeast1` |
| アクセス制御 | Uniform bucket-level access、Public access prevention = `enforced` |
| バージョニング | 有効（誤上書き・誤削除の復旧用）。非現行バージョンは 90 日で削除するライフサイクルルール |
| 管理 | `terraform/google/server/past-exams-bucket/{prd,stg,dev}` を新設（tfstate の prefix は既存に倣う） |

IAM:

| 主体 | ロール | 用途 |
| --- | --- | --- |
| `app-api-{env}` SA | `roles/storage.objectViewer`（バケット単位） | 署名付き GET URL の対象オブジェクト読み取り |
| `app-api-{env}` SA | `roles/iam.serviceAccountTokenCreator`（**自分自身の SA に対して**） | 鍵ファイルなしで `signBlob` による V4 署名を行うため |
| `admin-api-{env}` SA | `roles/storage.objectAdmin`（バケット単位） | 署名付き PUT URL の発行、削除 |
| `admin-api-{env}` SA | `roles/iam.serviceAccountTokenCreator`（自分自身） | 同上 |
| 移行作業者 | 一時的に `roles/storage.objectAdmin` | R2 からのコピー（移行完了後に削除） |

> 署名付き URL は SA の権限で検証されるため、プロジェクト全体ではなくバケット単位で付与する。

#### オブジェクトキー

```
past-exams/{past_exam_id}/{file_id}.pdf
```

- ID はいずれも UUID。表示名・年度などはキーに含めず DB に持つ（リネームでオブジェクトを動かさない）。
- `Content-Type: application/pdf`、`Content-Disposition` はダウンロード URL 発行時に `response-content-disposition` で元ファイル名を指定する。

### 2. DB スキーマ（server、Atlas マイグレーション追加）

1 科目に「過去問セット」が複数（年度・試験ごと）あり、1 セットに複数ファイル（問題・解答・追試等）がある、という 2 階層にする。

#### `past_exams`

| 列 | 型 | 備考 |
| --- | --- | --- |
| `id` | uuid PK | `gen_random_uuid()` |
| `subject_id` | uuid NOT NULL | FK → `subjects.id`（`ON DELETE CASCADE`）。インデックス |
| `year` | bigint NOT NULL | 実施年度 |
| `exam_type` | text NOT NULL | `midterm` / `final` / `quiz` / `makeup` / `other` |
| `title` | text NOT NULL | 表示名（例: `2024 期末`）。空なら年度・種別から生成して表示 |
| `note` | text NOT NULL DEFAULT `''` | 補足 |
| `created_at` / `updated_at` | timestamptz | `Common` |

- 並び順: `year DESC, exam_type, title`
- インデックス: `(subject_id, year)`

#### `past_exam_files`

| 列 | 型 | 備考 |
| --- | --- | --- |
| `id` | uuid PK | |
| `past_exam_id` | uuid NOT NULL | FK → `past_exams.id`（`ON DELETE CASCADE`）。インデックス |
| `kind` | text NOT NULL | `question` / `answer` / `other` |
| `file_name` | text NOT NULL | 元のファイル名（表示・ダウンロード名） |
| `object_key` | text NOT NULL UNIQUE | GCS のキー |
| `content_type` | text NOT NULL | 当面 `application/pdf` のみ許可 |
| `size_bytes` | bigint NOT NULL | |
| `sha256` | text NOT NULL | 重複検出・移行検証用 |
| `status` | text NOT NULL | `pending`（URL 発行済み・未アップロード） / `ready` |
| `display_order` | bigint NOT NULL DEFAULT 0 | |
| `created_at` / `updated_at` | timestamptz | |

- アプリには `status = 'ready'` のみ返す。
- 科目のバージョン（年度ごとの `subjects`）をまたいで同じ過去問を見せたい場合は、当面は登録時に科目ごとに紐付ける。将来 `lesson_id` 単位で共有する必要が出たら `past_exams.subject_id` を「科目群」への FK に変える（「未決事項」参照）。

GORM モデルは `internal/shared/model/past_exam.go` に `PastExam` / `PastExamFile` として追加し、`cmd/atlas-loader` に登録、`mise run migrate:diff` でマイグレーションを生成する。ER 図（`docs/er-diagram.md`）も更新する。

### 3. App API（Connect RPC）

`proto/app/v1/past_exam.proto` を追加する。

```proto
// 過去問を提供する。大学アカウントでログインしたユーザーのみ利用できる。
service PastExamService {
  // 科目の過去問を一覧する。
  rpc ListPastExams(ListPastExamsRequest) returns (ListPastExamsResponse);
  // 過去問ファイルのダウンロード URL を発行する。
  rpc GetPastExamFileDownloadUrl(GetPastExamFileDownloadUrlRequest) returns (GetPastExamFileDownloadUrlResponse);
}

message PastExamFile {
  string id = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  PastExamFileKind kind = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  string file_name = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  int64 size_bytes = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
}

message PastExam {
  string id = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  int32 year = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  PastExamType exam_type = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  string title = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
  string note = 5 [(google.api.field_behavior) = OUTPUT_ONLY];
  repeated PastExamFile files = 6 [(google.api.field_behavior) = OUTPUT_ONLY];
}

message ListPastExamsRequest {
  // 科目 ID。
  string subject_id = 1 [(google.api.field_behavior) = REQUIRED];
}
message ListPastExamsResponse {
  repeated PastExam past_exams = 1;
}

message GetPastExamFileDownloadUrlRequest {
  // 過去問ファイル ID。
  string file_id = 1 [(google.api.field_behavior) = REQUIRED];
}
message GetPastExamFileDownloadUrlResponse {
  // 署名付き URL。GET で取得する。
  string url = 1;
  // URL の有効期限。
  google.protobuf.Timestamp expire_time = 2;
}
```

- 列挙型 `PastExamType` / `PastExamFileKind` は既存の proto に倣い `*_UNSPECIFIED = 0` を置く。`buf lint` / `api-linter` を通す。
- 一覧には URL を含めない（一覧表示のたびに署名を発行しない。期限切れ URL がキャッシュに残らない）。
- `Subject` メッセージに `past_exam_count`（`int32`）を追加し、科目詳細で過去問タブの出し分けに使えるようにする（一覧 RPC を呼ぶ前に有無が分かる）。

#### 認可

- 2 RPC を `auth.userRequiredProcedures` に追加する（App Check + Firebase ID トークン必須）。
- さらにハンドラー（または `auth` の procedure 単位の設定）で **メールアドレスのドメインが `fun.ac.jp` で終わり、`email_verified = true`** であることを確認し、満たさなければ `PermissionDenied`。許可ドメインは環境変数 `PAST_EXAM_ALLOWED_EMAIL_DOMAINS`（カンマ区切り）で設定する。
- 存在しない ID は `NotFound`。`status != ready` のファイルも `NotFound` とする。

#### 署名付き URL

- `cloud.google.com/go/storage` の `BucketHandle.SignedURL`（`SigningScheme: SigningSchemeV4`、`Method: GET`）。Cloud Run 上では鍵ファイルを使わず、`GoogleAccessID` に実行 SA を指定し IAM Credentials API の `signBlob` で署名する（ライブラリが自動で行う）。
- 有効期限は **5 分**（`PAST_EXAM_SIGNED_URL_TTL`、既定 `5m`）。アプリは取得直後にダウンロードするので短くてよい。
- `response-content-disposition=attachment; filename*=UTF-8''<file_name>` と `response-content-type=application/pdf` を付ける。
- 発行時に `slog` で `user_id` / `file_id` を記録する（監査用。GCS のデータアクセス監査ログはコストを見て stg で試してから prd で有効化を判断）。

#### 実装配置

- `internal/modules/appapi/handler/past_exam.go`: ハンドラー（既存の `subject.go` と同様に GORM で直接参照）。
- `internal/shared/objectstorage/`（新設）: `Signer` インターフェース（`SignedGetURL(ctx, key, opts)` / `SignedPutURL(...)` / `Attrs(...)` / `Delete(...)`）と GCS 実装。テストはフェイクを差し込む。
- `cmd/app-api/main.go`: 環境変数 `PAST_EXAM_BUCKET` を読んで GCS クライアントを生成し、`Handler` に渡す。

### 4. Admin API（登録・管理）

`proto/admin/v1/past_exam.proto` に CRUD とアップロード用 RPC を追加する。admin は Cloudflare Access で保護済みのため認可は既存どおり。

| RPC | 内容 |
| --- | --- |
| `ListPastExams` / `GetPastExam` | 科目 ID 等で検索 |
| `CreatePastExam` / `UpdatePastExam` / `DeletePastExam` | メタデータの操作。削除時は配下のファイルの GCS オブジェクトも削除 |
| `CreatePastExamFileUpload` | `past_exam_files` を `pending` で作成し、署名付き PUT URL（15 分、`Content-Type: application/pdf` 固定）を返す |
| `CompletePastExamFileUpload` | GCS のオブジェクト属性（存在・サイズ・Content-Type）を確認し、`sha256` を計算して `ready` にする |
| `DeletePastExamFile` | DB 行と GCS オブジェクトを削除 |

- サイズ上限 50 MB（PUT URL に `x-goog-content-length-range` を付ける）。
- `pending` のまま 24 時間経過した行とオブジェクトは掃除する（既存の batch-jobs に追加するか、当面は手動。「未決事項」参照）。

### 5. アプリ（app リポジトリ）

| 変更 | 内容 |
| --- | --- |
| データ層 | `PastExamDataSource` を Connect クライアント（`PastExamService`）に置き換える。`S3DataSource` / `s3ClientProvider` / `minio` 依存を削除 |
| ドメイン | `PastExamRepository.fetchKeys(String)` → `fetchPastExams(subjectId)`（`PastExam` / `PastExamFile` エンティティを追加）と `fetchDownloadUrl(fileId)` |
| PDF | `PdfRepository.download` の `isPastExam` 分岐を削除。ダウンロード URL を取得してから既存の `_downloadWeb`（HTTP GET）で取得する。`CloudflarePdfViewer` は `PastExamPdfViewer` に改名し、URL 取得 → `WebPdfViewer` の順に呼ぶ |
| 紐付け | `SubjectDataSource.readPastExamId` と `Subject.pastExamId` を削除。科目 ID（UUID）で一覧を取得する。過去問タブは `pastExamCount > 0` で表示 |
| ルーティング | `:pastExamId`（オブジェクトキー）→ `:fileId`（UUID）。ディープリンクの互換は不要（キーは外部共有していない想定） |
| 画面 | 一覧を年度・試験種別でグルーピングし、ファイルごとに「問題 / 解答」を表示 |
| 認可エラー | `PermissionDenied` の場合は「大学の Google アカウントでログインしてください」を表示 |
| 設定 | `CLOUDFLARE_R2_*` の `--dart-define` と CI シークレットを削除 |
| テスト | `fake_past_exam_data_source.dart` を新インターフェースに合わせて更新 |

> 科目自体がまだ `syllabus.db`（`LessonId`）由来の画面では、API の科目 ID が取れない。アプリの科目取得が App API の `SubjectService` に移行済みであることを前提とする。未移行の画面が残る場合は、移行期間に限り `ListPastExamsRequest` に `lesson_id` + `year` での指定を追加する（`oneof`）。

### 6. データ移行（R2 → GCS）

一度きりのスクリプト `cmd/migrate-past-exams`（または `scripts/`）で行い、移行後に削除する。

1. R2 のバケットを列挙し、`{pastExamId}/{ファイル名}` のキー一覧を作る。
2. `syllabus.db` の `detail` から `過去問` → `LessonId` の対応を取り、サーバーの `syllabuses.lesson_id` → `subjects.id` を引く。同じ `lesson_id` に複数年度の科目がある場合は **最新年度の科目** に紐付ける。対応が取れないものはレポートに出して手動で判断する。
3. ファイル名から年度・試験種別・問題 / 解答を推定する（推定できなければ `year` = 不明として `0`、`exam_type = other`、`kind = other`。管理画面で後から修正）。
4. 各ファイルを R2 から取得し、`sha256` を計算して GCS（`past-exams/{past_exam_id}/{file_id}.pdf`）にアップロード、DB に `ready` で登録する。同一 `sha256` は重複として 1 つにまとめる。
5. 件数・合計サイズ・ハッシュを R2 側と突き合わせたレポートを出す。
6. dev → stg → prd の順に実行する。スクリプトは冪等にする（`object_key` 生成を R2 キーからの UUIDv5（`internal/shared/uuidv5`）にし、再実行で重複しない）。

### 7. ロールアウト

| 段階 | 内容 |
| --- | --- |
| 1 | infra: バケット・IAM・`app-api` / `admin-api` の環境変数（`PAST_EXAM_BUCKET` 等）を追加 |
| 2 | server: マイグレーション、Admin API、App API を実装し stg にデプロイ |
| 3 | stg でデータ移行し、アプリ（stg ビルド）で確認 |
| 4 | prd へマイグレーション・デプロイ・データ移行 |
| 5 | app: API 版をリリース。旧バージョンは R2 で引き続き動作 |
| 6 | 強制アップデート（最低バージョン）で旧バージョンを締め出した後、**R2 のアクセスキーを失効**し、R2 バケットを読み取り専用で一定期間保持してから削除。`syllabus.db` の `過去問` 列も次の DB 更新で削除 |

R2 のキー失効までは課題 1 が残るため、段階 6 までの期間は短くする。並行して、現行キーの権限が読み取り専用・当該バケット限定になっていることを確認する。

## 代替案

| 案 | 不採用の理由 |
| --- | --- |
| API サーバーがファイル本体を中継する | Cloud Run のレスポンスサイズ・転送コスト・インスタンス時間が増える。署名付き URL で十分 |
| GCS を公開にして推測困難なキーで配る | URL が漏れると無期限で取得でき、ログイン必須を担保できない |
| R2 のまま API だけ作る（R2 の署名付き URL） | 権限・監査・Terraform が Cloudflare と GCP に分かれたまま。他データと同じ GCP に寄せる方針に合わない |
| Firebase Storage + セキュリティルール | クライアントが直接読むためドメイン制限などのロジックがルールに分散する。サーバーの DB（科目）との整合をルールで書けない |
| Cloud CDN + 署名付き Cookie | アクセス量に対して過剰。必要になれば後から前段に置ける |

## 未決事項

- 過去問を年度ごとの `subjects` ではなく `lesson_id`（年度横断の科目）単位で共有するか。シラバス再設計で同一 `lesson_id` の科目が年度ごとに作られるため、毎年紐付け直しが必要になる。科目群テーブルの導入有無と合わせて決める。
- 閲覧ドメインを `fun.ac.jp` に限定するか（現状はログインのみで閲覧可能）。教員アカウントの扱い。
- `pending` ファイルの掃除をバッチジョブ化するか。
- アプリ内でのダウンロード済み PDF のキャッシュ（現状は毎回取得）。
- 学生からの過去問投稿機能（将来。Admin API のアップロードフローを App API に流用できる形にしておく）。
