# 過去問機能のサーバー経由化と GCS 移行(サーバー側設計)

- 日付: 2026-11-09
- 対象: App API (`app.v1`)、Admin API、DB、インフラ
- 関連: アプリ側設計 `fun-dotto/app` の `docs/design/20261109_past_exam_server_api.md`

## 1. 背景と目的

現状、過去問は Flutter アプリが Cloudflare R2 (S3 互換) へ直接アクセスして取得している。サーバーには過去問に関する API・テーブル・proto が存在しない。

| 課題 | 内容 |
|---|---|
| 認証情報の露出 | R2 のアクセスキーがビルド時の `--dart-define` でアプリに埋め込まれ、バイナリから抽出できる |
| 認可が UI のみ | 未ログイン時に一覧を出さないのは画面側の制御で、ストレージへのアクセスは認証と結びついていない |
| ストレージ依存 | R2 から GCS への移行を検討しており、クライアントが保存先に直接依存していると移行コストが高い |

目的は次の3点。

1. クライアントからストレージの認証情報をなくす。
2. 過去問の取得を認証済みユーザーに限定し、サーバー側で強制する。
3. 保存先(R2/GCS)をクライアントから隠蔽し、移行をアプリに影響させない。

## 2. 方針

- オブジェクトキーは外部に出さず、不透明な `past_exam_id`(UUID)で参照する。
- ファイル本体はサーバーで中継せず、短命の署名付き URL を返し、クライアントが直接ダウンロードする。
- 過去問のメタデータは DB で管理し、実行時にバケットの一覧取得は行わない。

## 3. 構成

```
App ──(1) ListPastExams(subject_id)───────▶ app-api ──▶ Postgres(past_exams)
App ──(2) GetPastExamDownloadUrl(id)──────▶ app-api ──▶ GCS V4 署名 (IAM signBlob)
App ──(3) GET 署名付き URL (TTL 5分) ─────────────────▶ GCS (非公開バケット)
```

### 署名付き URL を採用する理由

- Cloud Run のレスポンスサイズ上限と、PDF 中継による帯域・CPU 負荷を避けられる。
- 代替案(サーバーがバイト列を中継)は、サイズ上限とコストの問題から採用しない。

### 署名方法

- Cloud Run のサービスアカウントで V4 署名を行う。鍵ファイルは使わず、IAM Credentials の `signBlob` を使う。
- 署名用に、同サービスアカウント自身への `roles/iam.serviceAccountTokenCreator` を付与する。
- 読み取り用に、バケットへ `roles/storage.objectViewer` を付与する。
- `cloud.google.com/go/storage` は `go.mod` に間接依存として既にあるため、直接依存へ昇格する。

## 4. API(`proto/app/v1/past_exam.proto` を新設)

既存の規約(`OUTPUT_ONLY` の付与、日本語コメント、サービス単位のファイル)に合わせる。

```proto
// 過去問を提供する。
service PastExamService {
  // 科目の過去問一覧を取得する。
  rpc ListPastExams(ListPastExamsRequest) returns (ListPastExamsResponse);
  // 過去問をダウンロードするための署名付き URL を発行する。
  rpc GetPastExamDownloadUrl(GetPastExamDownloadUrlRequest) returns (GetPastExamDownloadUrlResponse);
}

message PastExam {
  string id = 1;                        // 不透明な ID (UUID)
  string subject_id = 2;
  string file_name = 3;                 // 表示名
  google.protobuf.Timestamp created_at = 4;
}

message GetPastExamDownloadUrlResponse {
  string url = 1;
  google.protobuf.Timestamp expires_at = 2;
}
```

- 両 RPC を `internal/modules/appapi/auth/auth.go` の `userRequiredProcedures` に追加する。現行の「fun.ac.jp アカウントでサインイン」要件をサーバーで強制する形になる。メールドメインの検証が必要かは、既存の Firebase ID トークンの内容を確認して決める。
- `ListPastExams` の `subject_id` は `Subject.id`(UUID)を使う。
- 署名付き URL の TTL は 5 分。署名時に `response-content-disposition` でファイル名を固定する。
- 表示名(`file_name`)は、現在クライアントの `pastExamFileName` が行っている「オブジェクトキーの最初の `/` より後ろ」の抽出をサーバーへ移し、`past_exams.file_name` に保持する。

## 5. データモデル

```
past_exams
  id           uuid PK
  subject_id   uuid FK → subjects
  object_name  text NOT NULL   -- バケット内のキー。外部には出さない
  file_name    text NOT NULL
  size_bytes   bigint
  created_at   timestamptz
  UNIQUE (object_name)
```

- モデルは `internal/shared/model` に追加し、Atlas でマイグレーションを作成する。`docs/er-diagram.md` も更新する。
- 代替案(実行時に GCS をプレフィックスで一覧取得)は、毎回 GCS を呼ぶこと、ID が不安定になること、DB で管理できないことから採用しない。

### 要確認のリスク

アプリの現行ロジックは、同梱 `syllabus.db` の `過去問` 列があればそれを、なければ `syllabus.id` をプレフィックスに使っている。リポジトリ内の `assets/syllabus.db` では 341 行すべてが NULL だった(同梱コピーの状態であり、本番のバケットの構造とは限らない)。R2 の実際のキー一覧で、`syllabus.id` 以外のプレフィックスが使われていないか確認する。

## 6. 実装構成

- `internal/modules/appapi/handler/past_exam.go`: `PastExamService` を実装する。`Handler` に `PastExamStorage` インターフェースを注入し、テストでは署名をモックする。
- `PastExamStorage`: `SignedDownloadURL(ctx, objectName, filename string, ttl time.Duration) (url string, expiresAt time.Time, err error)`。GCS 実装と、ローカル開発用のスタブ(または fake-gcs-server)を用意する。
- 署名付き URL 発行時に `user_id` と `past_exam_id` を構造化ログへ記録する。URL 自体はログに出さない。
- 追加した RPC に対応する handler テストと、`userRequiredProcedures` のテストを追加する。

## 7. R2→GCS 移行手順

1. 非公開 GCS バケットを作成する(Uniform bucket-level access、公開アクセス防止)。CORS は不要(ネイティブアプリが直接取得するため)。
2. R2→GCS を `rclone` または Storage Transfer Service で一括コピーし、件数・サイズを突合する。
3. バケットの一覧から `past_exams` を生成する移行スクリプトを実行する。キーのプレフィックスと `syllabus_id` の対応を取り、マッチしなかったものはレポートに出す。
4. サーバー側(proto・DB・ハンドラー)をリリースする。
5. アプリを新 API に切り替える。
6. 旧アプリの利用が十分減ったら、R2 の読み取りキーを失効させてバケットを廃止する。

旧バージョンのアプリは R2 を直接参照するため、切り替え後もしばらく R2 を残す。廃止時期は強制アップデートの判断と合わせて決める。

## 8. 実装順序

1. `past_exam.proto` を追加し、`buf` でコードを生成する。
2. `past_exams` のモデルとマイグレーションを追加する。
3. `PastExamStorage`(GCS 署名)と `PastExamService` ハンドラー、認可登録、テストを実装する。
4. 移行スクリプトでデータを投入する。
5. (別フェーズ)Admin API に過去問登録用 RPC を追加する。初期は移行スクリプトによる一括投入で足りる。

## 9. セキュリティ・運用

- 署名付き URL は TTL 内なら誰でも使えるため、ログに出さず、TTL を短くする。
- 取得ログ(`user_id`、`past_exam_id`)で、誰がどの過去問を取得したかを追えるようにする。
- 必要に応じて、レート制限や取得回数の上限を追加する。

## 10. 未決事項

- 本番バケットの実際のキー構造と、`syllabus.id` 以外のプレフィックスの有無
- Admin API 経由の登録を今回の範囲に含めるか
- 署名付き URL の TTL(5 分は初期案)と、旧アプリのサポート終了時期
