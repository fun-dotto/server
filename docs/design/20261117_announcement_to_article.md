# Announcement を廃止し、microCMS の Article に置き換える

## 背景

現状のおしらせ（Announcement）は、PostgreSQL の `announcements` テーブルを正本として、次の経路で提供している。

| 提供元 | 経路 | 利用者 |
| --- | --- | --- |
| App API（Connect） | `app.v1.AnnouncementService/ListAnnouncements` → DB 直接参照 | Web（`/setting/announcements`）、モバイルアプリ（未リリース） |
| App BFF API（OpenAPI） | `GET /v1/announcements` → Announcement API（HTTP） | モバイルアプリ |
| Announcement API（OpenAPI） | CRUD → DB | App BFF API |
| Admin API（Connect） | `admin.v1.AnnouncementService` の CRUD → DB | admin-web |

App API を使うモバイルアプリはまだリリースされていないため、App API は後方互換性を保たずに置き換える。App BFF API はリリース済みのアプリが使うため、後方互換性を保つ。

おしらせは外部 URL へのリンクを並べるだけで、本文を Dotto 側で持てない。また、入稿のために Admin API と Announcement API の 2 系統を保守している。

## 基本方針

- おしらせの正本を **microCMS の `articles` API** に移し、DB の `announcements` は廃止する。
- 記事本文は Web でレンダリングする。記事の URL は `https://{dev.,stg.,}fun-dotto.com/articles/<id>` とする。
- **Web**: `/articles` で一覧、`/articles/<id>` で記事を表示する。Server-side から microCMS に直接リクエストし、Dotto の API は介さない。
- **App API**: `AnnouncementService` を削除し、`ArticleService` を追加する（アプリ用）。
- **App BFF API**: OpenAPI スキーマとドメインモデルは変えず、リポジトリ層の取得元を microCMS に差し替える（後方互換性を担保）。
- **Announcement API** と **Admin API の Announcement** を削除する。
- **GORM モデル** `model.Announcement` を削除し、`announcements` テーブルを DROP する。
- 既存の Announcement データは移行せず破棄する。Web の旧ページ `/setting/announcements` は削除する。

## 変更点

### 1. microCMS

#### API `articles`（リスト形式、新規作成）

| フィールド ID | 種類 | 必須 | 備考 |
| --- | --- | --- | --- |
| `title` | テキストフィールド | ○ | タイトル |
| `body` | リッチエディタ | ○ | 本文（Web でのみ使用） |

- `id` / `publishedAt` / `createdAt` / `updatedAt` / `revisedAt` は microCMS の組み込みフィールドを使う。
- 公開日時は `publishedAt`（初回公開日時）を使う。予約公開は microCMS の予約公開機能で行う。
- 公開終了（旧 `available_until`）は microCMS の公開終了予約（非公開化）で代替する。
- サービスは dev / stg / prd で **1 つを共有する**。記事 URL のホストだけを環境ごとの `ARTICLE_BASE_URL` で切り替える。
  - 記事は読み取り専用で全環境に同じものを見せるため、当面は共有で問題ない。dev で確認のために公開した記事は本番にも表示される点に注意する（詳細ページは draftKey で公開前にプレビューできるが、一覧は公開しないと確認できない）。

#### 認証情報

| 環境変数 | 内容 | 利用者 |
| --- | --- | --- |
| `MICROCMS_SERVICE_DOMAIN` | サービスドメイン（`https://<domain>.microcms.io`） | app-api, app-bff-api, web |
| `MICROCMS_API_KEY` | `articles` の GET 権限のみを持つ API キー | app-api, app-bff-api, web |
| `ARTICLE_BASE_URL` | 記事 URL のベース（例: `https://dev.fun-dotto.com`） | app-api, app-bff-api |

- API キーは Secret Manager に保存し、Cloud Run（server）と Web のデプロイ先に注入する（infra リポジトリで管理）。
- 取得は `GET /api/v1/articles` で `fields=id,title,publishedAt`、`orders=-publishedAt` を指定する。下書き・非公開の記事は返らない。
- `limit` の上限が 100 のため、全件取得は `offset` を進めて `totalCount` に達するまで繰り返す。

### 2. 共通 microCMS クライアント（`internal/shared/microcms`、新規）

App API と App BFF API で共有する薄いクライアントを置く。

```go
package microcms

type Article struct {
	ID          string
	Title       string
	PublishedAt time.Time
}

type Client struct { /* serviceDomain, apiKey, httpClient */ }

func NewClientFromEnv() (*Client, error)

// ListArticles は公開中の記事を publishedAt の降順で全件返す。
func (c *Client) ListArticles(ctx context.Context) ([]Article, error)
```

- リクエストヘッダー `X-MICROCMS-API-KEY` に API キーを付ける。
- タイムアウトを設定し、2xx 以外はエラーにする。
- 記事 URL の組み立て（`ARTICLE_BASE_URL + "/articles/" + id`）は `func ArticleURL(baseURL, id string) string` として同パッケージに置く。
- テストは `httptest.Server` で microCMS のレスポンスを模擬する。

### 3. App API（Connect）

#### proto

`proto/app/v1/announcement.proto` を削除し、`proto/app/v1/article.proto` を追加する。

```proto
syntax = "proto3";

package app.v1;

import "google/api/field_behavior.proto";
import "google/protobuf/timestamp.proto";

// 記事を提供する。
service ArticleService {
  // 公開中の記事の一覧を、公開日時の新しい順に取得する。
  rpc ListArticles(ListArticlesRequest) returns (ListArticlesResponse);
}

// 記事。
message Article {
  // ID。microCMS のコンテンツ ID。
  string id = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  // 記事ページの URL。
  string url = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  // タイトル。
  string title = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  // 公開日時。
  google.protobuf.Timestamp publish_time = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
}

// ListArticles の Request。
message ListArticlesRequest {
}

// ListArticles の Response。
message ListArticlesResponse {
  // 記事の一覧。
  repeated Article articles = 1;
}
```

- フィールドは `id` / `url` / `title` / `publish_time` とする。`id` は microCMS のコンテンツ ID をそのまま返す（アプリの詳細画面が `id` で記事を探すため）。公開日時は AIP-142 に合わせて `publish_time` とする。
- ページングは設けない（app-api-design.md の方針どおり）。

#### 実装

- `internal/modules/appapi/handler/announcement.go` を削除し、`article.go` を追加する。
- `Handler` に `microcms.Client`（interface で受ける）と `articleBaseURL` を持たせる。
- `cmd/app-api/main.go`: `NewAnnouncementServiceHandler` を `NewArticleServiceHandler` に置き換える。
- `internal/modules/appapi/auth/auth_test.go` のダミーハンドラーを Article に置き換える。
- 認証は既存のサービスと同じ Interceptor を適用する。

### 4. App BFF API（OpenAPI）

OpenAPI スキーマ（`api/openapi/app/openapi.yaml` の `GET /v1/announcements` と `Announcement`）とドメインモデル（`domain.Announcement`）は **変更しない**。

#### リポジトリ層

`internal/modules/app/repository/announcement.go` を、Announcement API クライアントではなく microCMS クライアントで実装し直す。

| `domain.Announcement` | 取得元 |
| --- | --- |
| `ID` | 記事 URL から生成した UUIDv5（`uuidv5.NewSHA1(uuidv5.NamespaceURL, []byte(url))`） |
| `Title` | `title` |
| `AvailableFrom` | `publishedAt` |
| `AvailableUntil` | 常に `nil` |
| `URL` | `ARTICLE_BASE_URL + "/articles/" + id` |

- `GetAnnouncements(query)`: microCMS は公開中の記事のみ返すため `FilterIsActive` は常に満たされる。`SortByDate` は `orders` パラメーター（`publishedAt` / `-publishedAt`）に変換する。
- `GetAnnouncement` / `CreateAnnouncement` / `UpdateAnnouncement` / `DeleteAnnouncement` は BFF のハンドラーから呼ばれていないため、`AnnouncementRepository` interface と `AnnouncementService` から削除する（`domain.AnnouncementRequest` と `external/announcement_model.go` も削除）。
- `announcement_mock.go` を再生成する。

#### レスポンスの後方互換性

| 項目 | 変更前 | 変更後 |
| --- | --- | --- |
| `id` | DB の UUID | 記事 URL から生成した UUIDv5（同じ記事なら常に同じ値） |
| `url` | 任意の外部 URL | `https://{env}fun-dotto.com/articles/<id>` |
| `date` | `available_from` | `publishedAt` |

スキーマ上の型（`string` / `uri` / `date-time`）と `id` の UUID 形式は変わらない。`id` は記事 URL から決定的に生成するため、リクエストごとに変わらない。

#### 起動処理

- `cmd/app-bff-api/main.go`: `clients.Announcement` の代わりに `microcms.NewClientFromEnv()` を渡す。
- `internal/shared/apiclient/clients.go` から Announcement API クライアントと `ANNOUNCEMENT_API_URL` を削除する。

### 5. Announcement API の削除

次を削除する。

- `cmd/announcement-api/`
- `internal/modules/announcement/`
- `api/openapi/announcement/`、`gen/announcement/`
- コード生成設定（`config.yaml` を参照している生成スクリプト・mise タスク）
- `internal/shared/server/server.go` のコメント中の `announcement-api`
- Dockerfile・CI のビルド対象

### 6. Admin API

- `proto/admin/v1/announcement.proto` と生成コードを削除する。
- `internal/modules/adminapi/handler/announcement.go` を削除する。
- `cmd/admin-api/main.go` から `NewAnnouncementServiceHandler` の登録を削除する。
- 入稿は microCMS の管理画面で行う。
- admin-web の Announcement 画面を先に削除してから、Admin API の Announcement を削除する。

### proto の破壊的変更

`app/v1/announcement.proto` と `admin/v1/announcement.proto` の削除は `buf breaking`（FILE）に違反する。次の理由から例外として扱う。

- `app.v1.AnnouncementService` を使うのは Web と未リリースのモバイルアプリだけで、どちらも削除前に `ArticleService` へ移行する。
- `admin.v1.AnnouncementService` を使うのは admin-web だけで、削除前に画面を削除する。

`buf breaking` は CI で実行していないため、CI への対応は不要。手元で実行した場合の失敗は、この 2 ファイルの削除によるものに限り許容する。

### 7. GORM モデル・DB

- `internal/shared/model/announcement.go` を削除する。
- `cmd/atlas-loader/main.go` から `&model.Announcement{}` を削除する。
- Atlas でマイグレーションを生成する（`DROP TABLE "public"."announcements";`）。データは移行せず破棄する。
- `docs/er-diagram.md` を更新する。

### 8. Web（`fun-dotto/web`）

- `src/lib/microcms.ts` を追加し、Server-side（Server Component）からのみ microCMS を呼ぶ。API キーをクライアントバンドルに含めない（`import "server-only"`）。
- `/articles`（`src/app/articles/page.tsx`）: 一覧。`fields=id,title,publishedAt`、`orders=-publishedAt`。
- `/articles/[id]`（`src/app/articles/[id]/page.tsx`）: `GET /api/v1/articles/{id}` で `title` / `body` / `publishedAt` を取得してレンダリングする。存在しない ID は `notFound()`。
- リッチエディタの HTML はサニタイズしてから描画する。
- キャッシュは ISR（`revalidate`）で行い、必要なら microCMS の Webhook から On-demand Revalidation を叩く。
- `generateMetadata` で記事タイトルを `<title>` / OGP に設定し、`sitemap.ts` に記事を追加する。
- `/setting/announcements` と、`nav-sections.ts` の「お知らせ」を削除し、`src/api/app/v1/announcement_pb.*` を削除する（ナビに `/articles` を追加するかは別途判断）。
- アプリからは記事 URL を WebView / ブラウザで開くため、`/articles/<id>` はログインなしで閲覧できるようにする。

### 9. Infra（`fun-dotto/infra`）

- `terraform/google/server/announcement-service/{dev,stg,prd}` を削除する（Cloud Run・サービスアカウント等）。`digger.yml` から該当プロジェクトを削除する。
- `app-bff-api` の `remote-states.tf` / `main.tf` から announcement-service の参照（`ANNOUNCEMENT_API_URL`、Invoker 権限）を削除する。
- `MICROCMS_SERVICE_DOMAIN` / `MICROCMS_API_KEY`（Secret Manager）/ `ARTICLE_BASE_URL` を app-api・app-bff-api に追加する。Web のデプロイ先にも同じ値を設定する。
- `foundation/monitoring/alert-policies.tf`、`foundation/cloud-sql/locals.tf`、`docs/postgresql_roles.md` の announcement 関連を削除する。

## リリース手順

依存関係があるため、次の順で進める。

1. infra: microCMS の設定を追加し（fun-dotto/infra の microCMS 設定の PR）、`microcms-api-key` に API キーを登録する。必要な記事を microCMS に入稿しておく。
2. web: `/articles` と `/articles/<id>` を追加し、`/setting/announcements` を削除する（`app.v1.AnnouncementService` を呼ばなくする）。
3. app: `app.v1.AnnouncementService` から `app.v1.ArticleService` に移行する（未リリースのため、次のリリースまでに済ませる）。
4. server: microCMS クライアント、App API の `ArticleService`（`AnnouncementService` の削除を含む）、App BFF のリポジトリ差し替えをデプロイする。ここで BFF の取得元が microCMS に切り替わる。
5. admin-web: Announcement 画面を削除する。
6. server: Announcement API・Admin API の Announcement・GORM モデルを削除し、`announcements` を DROP する。
7. infra: announcement-service の各 root で `terraform destroy` を実行してから、announcement-service の Terraform 定義と参照を削除する。

4 の時点で既存のおしらせはアプリに表示されなくなる（データは移行しない）。2 と 3 が 4 より後になると、Web の旧ページと開発中のアプリが `AnnouncementService` を呼べなくなる。

## テスト

- `internal/shared/microcms`: 正常系、ページング（101 件以上）、非 2xx、タイムアウト。
- App API `ListArticles`: URL の組み立て、並び順。
- App BFF のリポジトリ: `ID` が UUID 形式で、同じ記事なら同じ値になること。
- App API `ListArticles`: `id` に microCMS のコンテンツ ID が入ること。
- App BFF `AnnouncementsV1List`: microCMS のモックから既存スキーマどおりのレスポンスになること（既存の `announcement_v1_test.go` / `service/announcement_test.go` を書き換える）。
- Web: 一覧・詳細・404 の表示を確認する。
