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

- **過去問のメタデータを Cloud SQL（`past_exam_files` テーブル）に持ち、試験が実施された年度の科目（`subjects`）に紐付ける。** `syllabus.db` の `過去問` 列は使わない。
  - [科目テーブル・シラバステーブル再設計](./20261109_subject_syllabus_redesign.md) により、`subjects` は `(lesson_id, year)` ごとに作られる。2024 年度の期末試験は 2024 年度の科目に紐付ける。毎年の紐付け直しは不要（過去の科目に紐付いたまま）。
  - 閲覧時は、紐付いたシラバスの `lesson_id` が同じ科目（= 同じ授業の他年度）の過去問もまとめて返す。担当教員の違いなどは科目ごとの情報で分かる。
- **ファイルは GCS の非公開バケットに置く。** アプリは GCS に直接の資格情報を持たない。
- **App API（Connect RPC）に一覧取得とダウンロード URL 発行の RPC を追加する。** ダウンロードは **V4 署名付き URL（短期）** を返し、アプリはそれを HTTP GET する。ファイル本体を API サーバー経由で中継しない。
- **認可はサーバーで行う。** 過去問 RPC はログイン必須（`userRequiredProcedures`）にし、さらに学生アカウント（`b<数字7桁>@fun.ac.jp` / `g<数字7桁>@fun.ac.jp`）のユーザーに限る。
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
past-exams/{file_id}.pdf
```

- `file_id` は UUID（`past_exam_files.id`）。表示名・年度などはキーに含めず DB に持つ（リネームでオブジェクトを動かさない）。
- `Content-Type: application/pdf` で保存する。ダウンロード時のファイル名は URL 発行時に `response-content-disposition` で `title` から指定する。

### 2. DB スキーマ（server、Atlas マイグレーション追加）

過去問は構造化されたデータではないため、**1 科目に複数の過去問ファイルが所属する 1 階層** にする。テーブルは `past_exam_files` の 1 つだけ。

#### `past_exam_files`

| 列 | 型 | 備考 |
| --- | --- | --- |
| `id` | uuid PK | `gen_random_uuid()` |
| `subject_id` | uuid NOT NULL | FK → `subjects.id`（`ON DELETE CASCADE`）。インデックス。実施年度は `subjects.year` で決まる |
| `title` | text NOT NULL | 表示名（例: `2024 期末 問題`）。ダウンロード時のファイル名にも使う |
| `object_key` | text NOT NULL UNIQUE | GCS のキー。Content-Type などのオブジェクト属性が必要になれば、ここから GCS を引く |
| `size_bytes` | bigint NOT NULL | |
| `sha256` | text NOT NULL | 重複検出・移行検証用 |
| `created_at` / `updated_at` | timestamptz | `Common` |

- 並び順: `subjects.year DESC, title`
- 行は GCS へのアップロード完了後に作成するため、アップロード状態の列は持たない（「Admin API」参照）。
- 同じ授業の他年度の過去問は、`subjects` → `syllabuses`（`syllabus_id`）を JOIN し、`syllabuses.lesson_id` が同じ科目の `past_exam_files` として引く。`subjects.syllabus_id` は UNIQUE で、付け替えも同じ `(lesson_id, year)` 内に限られる（再設計の admin-api バリデーション）ため、シラバスのバージョンが変わっても紐付けは崩れない。

```sql
SELECT f.*, s.year
FROM past_exam_files f
JOIN subjects s ON s.id = f.subject_id
JOIN syllabuses sy ON sy.id = s.syllabus_id
WHERE sy.lesson_id = (
  SELECT sy2.lesson_id FROM subjects s2 JOIN syllabuses sy2 ON sy2.id = s2.syllabus_id WHERE s2.id = $1
)
ORDER BY s.year DESC, f.title;
```

GORM モデルは `internal/shared/model/past_exam_file.go` に `PastExamFile` として追加し、`cmd/atlas-loader` に登録、`mise run migrate:diff` でマイグレーションを生成する。ER 図（`docs/er-diagram.md`）も更新する。

### 3. App API（Connect RPC）

`proto/app/v1/past_exam.proto` を追加する。

```proto
// 過去問を提供する。大学アカウントでログインしたユーザーのみ利用できる。
service PastExamService {
  // 科目の過去問を一覧する。
  rpc ListPastExamFiles(ListPastExamFilesRequest) returns (ListPastExamFilesResponse);
  // 過去問ファイルのダウンロード URL を発行する。
  rpc GetPastExamFileDownloadUrl(GetPastExamFileDownloadUrlRequest) returns (GetPastExamFileDownloadUrlResponse);
}

// 過去問ファイル。
message PastExamFile {
  string id = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  // 紐付いた科目の ID。
  string subject_id = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  // 実施年度（紐付いた科目の年度）。
  int32 year = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  // 表示名。
  string title = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
  int64 size_bytes = 5 [(google.api.field_behavior) = OUTPUT_ONLY];
}

message ListPastExamFilesRequest {
  // 対象の授業。どちらか一方を指定する。
  // 同じ授業（lesson_id）の全年度の科目に紐付いた過去問を返す。
  oneof target {
    // 科目 ID。
    string subject_id = 1;
    // 授業コード（教務システムの lesson_id）。syllabus.db 由来の科目画面からの移行期間用。
    string lesson_id = 2;
  }
}
message ListPastExamFilesResponse {
  repeated PastExamFile past_exam_files = 1;
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

- `buf lint` / `api-linter` を通す。
- 一覧には URL を含めない（一覧表示のたびに署名を発行しない。期限切れ URL がキャッシュに残らない）。
- `Subject` メッセージに `past_exam_file_count`（`int32`）を追加し、科目詳細で過去問タブの出し分けに使えるようにする（一覧 RPC を呼ぶ前に有無が分かる）。

#### 認可

- 2 RPC を `auth.userRequiredProcedures` に追加する（App Check + Firebase ID トークン必須）。
- さらにハンドラー（または `auth` の procedure 単位の設定）で、Firebase ID トークンの **`email` が `^[bg][0-9]{7}@fun\.ac\.jp$` に一致し、`email_verified = true`** であることを確認し、満たさなければ `PermissionDenied`。
  - `b` は学部生、`g` は大学院生のアカウント。教員・職員など、これ以外の `@fun.ac.jp` アカウントは閲覧できない。
  - 大文字小文字の揺れに備え、比較前に `email` を小文字にする。
  - `auth.User` に `EmailVerified` を追加し、`FirebaseIDTokenVerifier` でクレームから詰める。判定は `IsStudent()` などの関数に切り出してテストする。
- 存在しない ID は `NotFound`。

#### 署名付き URL

- `cloud.google.com/go/storage` の `BucketHandle.SignedURL`（`SigningScheme: SigningSchemeV4`、`Method: GET`）。Cloud Run 上では鍵ファイルを使わず、`GoogleAccessID` に実行 SA を指定し IAM Credentials API の `signBlob` で署名する（ライブラリが自動で行う）。
- 有効期限は **5 分**（`PAST_EXAM_SIGNED_URL_TTL`、既定 `5m`）。アプリは取得直後にダウンロードするので短くてよい。
- `response-content-disposition=attachment; filename*=UTF-8''<title>.pdf` と `response-content-type=application/pdf` を付ける。
- 発行時に `slog` で `user_id` / `file_id` を記録する（監査用。GCS のデータアクセス監査ログはコストを見て stg で試してから prd で有効化を判断）。

#### 実装配置

- `internal/modules/appapi/handler/past_exam.go`: ハンドラー（既存の `subject.go` と同様に GORM で直接参照）。
- `internal/shared/objectstorage/`（新設）: `Signer` インターフェース（`SignedGetURL(ctx, key, opts)` / `SignedPutURL(...)` / `Attrs(...)` / `Delete(...)`）と GCS 実装。テストはフェイクを差し込む。
- `cmd/app-api/main.go`: 環境変数 `PAST_EXAM_BUCKET` を読んで GCS クライアントを生成し、`Handler` に渡す。

### 4. Admin API（登録・管理）

`proto/admin/v1/past_exam.proto` に CRUD とアップロード用 RPC を追加する。admin は Cloudflare Access で保護済みのため認可は既存どおり。

| RPC | 内容 |
| --- | --- |
| `ListPastExamFiles` / `GetPastExamFile` | 科目 ID 等で検索 |
| `CreatePastExamFileUploadUrl` | 新しい `file_id` を採番し、`past-exams/{file_id}.pdf` への署名付き PUT URL（15 分、`Content-Type: application/pdf` 固定）を返す。DB には書かない |
| `CreatePastExamFile` | `file_id`・`subject_id`・`title` を受け取り、GCS にオブジェクトが存在することを確認して `size_bytes`・`sha256` を計算し、行を作成する |
| `UpdatePastExamFile` | `title`・`subject_id` の変更 |
| `DeletePastExamFile` | DB 行と GCS オブジェクトを削除 |

- サイズ上限 50 MB（PUT URL に `x-goog-content-length-range` を付ける）。
- PUT 後に `CreatePastExamFile` が呼ばれなかったオブジェクトは DB から参照されないまま残る。バケットのライフサイクルルールでは判別できないため、DB に行のないオブジェクトを当面は手動で掃除する（「未決事項」参照）。

### 5. アプリ（app リポジトリ）

| 変更 | 内容 |
| --- | --- |
| データ層 | `PastExamDataSource` を Connect クライアント（`PastExamService`）に置き換える。`S3DataSource` / `s3ClientProvider` / `minio` 依存を削除 |
| ドメイン | `PastExamRepository.fetchKeys(String)` → `fetchPastExamFiles(subjectId)`（`PastExamFile` エンティティを追加）と `fetchDownloadUrl(fileId)` |
| PDF | `PdfRepository.download` の `isPastExam` 分岐を削除。ダウンロード URL を取得してから既存の `_downloadWeb`（HTTP GET）で取得する。`CloudflarePdfViewer` は `PastExamPdfViewer` に改名し、URL 取得 → `WebPdfViewer` の順に呼ぶ |
| 紐付け | `SubjectDataSource.readPastExamId` と `Subject.pastExamId` を削除。API の科目なら `subject_id`、`syllabus.db` 由来の科目なら `LessonId` を `lesson_id` に指定して一覧を取得する。API の科目では過去問タブを `pastExamFileCount > 0`（同じ授業の全年度分の件数）で表示 |
| ルーティング | `:pastExamId`（オブジェクトキー）→ `:fileId`（UUID）。ディープリンクの互換は不要（キーは外部共有していない想定） |
| 画面 | 一覧を年度でグルーピングし、`title` を表示 |
| 認可エラー | `PermissionDenied` の場合は「学生の大学アカウント（b/g + 学籍番号）でログインしてください」を表示 |
| 設定 | `CLOUDFLARE_R2_*` の `--dart-define` と CI シークレットを削除 |
| テスト | `fake_past_exam_data_source.dart` を新インターフェースに合わせて更新 |

> 過去問は Subject に紐付けるが、一覧は `lesson_id` でも引けるので、科目画面が App API の `SubjectService` に移行済みかどうかに関係なく過去問を API 化できる。`syllabus.db` の `LessonId` は教務システムの授業コードで、再設計後の `syllabuses.lesson_id` と同じ値になる。科目画面の移行が終わったら `lesson_id` 指定を削除し、`subject_id` に揃える。

### 6. データ移行（R2 → GCS）

一度きりのスクリプト `cmd/migrate-past-exams`（または `scripts/`）で行い、移行後に削除する。

1. R2 のバケットを列挙し、`{pastExamId}/{ファイル名}` のキー一覧を作る。
2. `syllabus.db` の `detail` から `過去問` → `LessonId` の対応を取る。ファイル名から実施年度を推定し、`syllabuses.lesson_id` と `subjects.year` が一致する科目に紐付ける。
   - 年度が推定できない、またはその年度の科目がない（シラバス取得ジョブで過去年度を取得していない等）ものはレポートに出す。必要な年度のシラバス取得・科目生成ジョブを実行してから再実行するか、手動で判断する。
3. `title` は R2 のファイル名（拡張子なし）をそのまま使う（管理画面で後から修正できる）。
4. 各ファイルを R2 から取得し、`sha256` を計算して GCS（`past-exams/{file_id}.pdf`）にアップロードしてから DB に登録する。同一 `sha256` は重複として 1 つにまとめる。
5. 件数・合計サイズ・ハッシュを R2 側と突き合わせたレポートを出す。
6. dev → stg → prd の順に実行する。スクリプトは冪等にする（`file_id` を R2 キーからの UUIDv5（`internal/shared/uuidv5`）にし、再実行で重複しない）。

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

- DB に行のない GCS オブジェクト（アップロード後に登録されなかったもの）の掃除をバッチジョブ化するか。
- アプリ内でのダウンロード済み PDF のキャッシュ（現状は毎回取得）。
- 学生からの過去問投稿機能（将来。Admin API のアップロードフローを App API に流用できる形にしておく）。
