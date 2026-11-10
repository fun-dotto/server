# Atlas マイグレーションの CI/CD 整備

## 背景

マイグレーションの生成は Atlas CLI、適用は自作バイナリ（`cmd/migrate`）と役割が分かれており、適用前に内容を確認する手段がない。
また CI でマイグレーションを一切検証しておらず、環境（dev / stg / prd）ごとの適用フローも定義されていない。

## 要件

| 環境 | トリガー | 動作 |
| --- | --- | --- |
| dev / stg / prd | `pull_request`, `merge_group` | dry-run（適用予定の SQL を表示するのみ。DB は変更しない） |
| dev / stg | `push`（main） | 適用 |
| prd | `workflow_dispatch` | 適用 |

## 現状

### 構成

| 役割 | 実装 |
| --- | --- |
| スキーマの源泉 | GORM モデル（`internal/shared/model`）を `cmd/atlas-loader` が列挙し、`atlas-provider-gorm` で SQL 化 |
| diff 生成 | `mise run migrate:diff`（`atlas migrate diff --env gorm`）。dev DB は `docker://postgres/18/dev` |
| 適用 | `mise run migrate:apply`（`go run ./cmd/migrate`）。Atlas CLI は使わず、Atlas の Go API（`migrate.NewExecutor`）で適用 |
| 本番適用 | 同一イメージの `/bin/migrate` を Cloud Run Job で実行。cloudsqlconn + IAM 認証で接続 |
| リビジョン管理 | `internal/migrate/revisions.go` に Atlas CLI 互換の `atlas_schema_revisions` を自前実装（約 200 行） |
| 権限 | 別ジョブ `apply-table-privileges` が全テーブルの owner を `dotto_admin` に変更し、権限を再付与 |
| `atlas.hcl` | `env "gorm"` のみ（`src` + `dev` + `migration.dir`）。接続先 `url` を持つ env はない |

### CI / CD

- `ci.yml`: test / build / lint のみ。マイグレーションの検証はない。
- `deploy.yml`: main への push で stg にデプロイ。`workflow_dispatch` で環境を選択でき、prd は `action: deploy` のみ指定可。**`migrate` / `apply-table-privileges` の Job を実行するステップはない**（`set-env` の `DEPLOY_TARGETS` にも含まれない）。適用は手動、または別管理と思われる。
- dev へは `workflow_dispatch` でのみデプロイ可能。push 時の対象は stg のみ。

### 課題

1. **dry-run できない。** `cmd/migrate` は常に適用する。`atlas.hcl` に接続先を持つ env がないため `atlas migrate apply --dry-run` も使えない。
2. **PR で検証されない。** `atlas migrate lint` / `validate`（`atlas.sum` 整合性）や、モデルとマイグレーションの差分残りの検査がない。生成し忘れや破壊的変更に PR で気付けない。
3. **適用が CI/CD に組み込まれていない。** 環境ごとの適用トリガーが定義されておらず、実行順序（マイグレーション → 権限付与 → アプリのデプロイ）の保証もない。
4. **Atlas CLI の再実装を抱えている。** リビジョンテーブルの読み書きを自前実装しており、Atlas のバージョンアップで互換性が崩れるリスクがある。ロックもなく、Job の多重起動に弱い。
5. Atlas の CLI（mise: 1.3.0）と Go ライブラリ（`ariga.io/atlas`）の二重管理。
6. マイグレーションファイル名が不揃い（`20260510085829.sql` のように名前なしのものがある）。

## 理想形

### 方針

- **適用も Atlas CLI に一本化する。** `atlas migrate apply` を使い、`--dry-run` で適用予定を確認できるようにする。自作のリビジョン管理（`revisions.go`）と `cmd/migrate` は廃止する。
- **環境ごとに `atlas.hcl` の env を定義する。** 接続先は環境変数（`DATABASE_URL`）で注入する。
- **GitHub Actions から Cloud SQL に IAM 認証で接続する。** Workload Identity Federation で認証し、Cloud SQL Auth Proxy 経由で接続する。環境ごとの GitHub Environment の `vars` で接続先を切り替える。
- **PR で静的検証と dry-run を行う。** DB を変更しない検証をマージ前に必ず通す。

### ワークフロー

| トリガー | 環境 | 実行内容 |
| --- | --- | --- |
| `pull_request` / `merge_group` | dev, stg, prd（matrix） | `atlas migrate validate` → `atlas migrate lint` → モデル差分チェック → `atlas migrate apply --dry-run` |
| `push`（main） | dev, stg（matrix） | `atlas migrate apply` → `apply-table-privileges` |
| `workflow_dispatch` | 入力で選択（prd 用） | `atlas migrate apply` → `apply-table-privileges`（GitHub Environment `prd` の承認を必須にする） |

補足:

- dry-run は DB に書き込まない。ただし接続して未適用のリビジョンを読むため、PR の CI にも各環境への読み取り接続が必要になる。fork からの PR ではシークレットが渡らないため、dry-run は同一リポジトリの PR のみ対象とする。
- 静的検証（validate / lint / 差分チェック）は環境に依存しないため 1 回だけ実行し、環境ごとの matrix は `apply --dry-run` のみに絞る。
- `merge_group` でも同じ検証を実行し、マージキュー上の最終状態を検証する。
- `workflow_dispatch` の prd は、dev / stg で適用済みのコミットのみ許可する運用を想定する（`deploy.yml` の prd 制約と揃える）。
- アプリのデプロイ（`deploy.yml`）より前にマイグレーションが完了するよう、同一ワークフロー内で `needs` により順序を保証する。

## 変更点

### 1. `atlas.hcl`

環境ごとの env を追加する。接続先は `getenv` で注入する。

```hcl
variable "database_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "dev" {
  url = var.database_url
  migration {
    dir = "file://migrations"
  }
}

# stg / prd も同様（接続先は GitHub Environment ごとの vars / secrets で切り替える）
```

`env "gorm"`（diff 生成用）は据え置く。lint 用に `lint { latest = 1 }` などの設定も追加する。

### 2. ワークフロー

- 新規 `.github/workflows/migrate.yml` を追加する。`deploy.yml` とは別に管理し、トリガーと環境の組み合わせを上表のとおり定義する。
- 共通処理は composite action（例: `.github/actions/atlas-migrate`）に切り出す。内容は次のとおり。
  1. WIF で GCP 認証
  2. Cloud SQL Auth Proxy 起動（IAM 認証）
  3. mise で Atlas をセットアップ（`mise.toml` のバージョンと同一）
  4. `atlas migrate apply --env <env>`（dry-run 時は `--dry-run`）
- PR 向けの静的検証ジョブでは、dev DB として `docker://postgres/18/dev` を使う（`atlas.hcl` の既存設定）。
- `concurrency` を環境ごとに設定し、同一環境へのマイグレーションの並列実行を防ぐ（Atlas は advisory lock も取得する）。
- `deploy.yml` の push 対象を stg から dev + stg に拡張するかは別途判断する。マイグレーションのみ dev / stg 両方を対象とし、アプリのデプロイ範囲は現状維持でもよい。

### 3. 権限付与

`apply-table-privileges` は、移行期間中は現状のバイナリを流用してマイグレーション適用後に実行する。最終的には、マイグレーションの SQL に `GRANT` / `ALTER OWNER` を含める、または `ALTER DEFAULT PRIVILEGES` に寄せて別ジョブを不要にすることを検討する。

### 4. コードの削除

- `cmd/migrate`、`internal/migrate/{apply,revisions}.go` を削除する。
- `go.mod` から `ariga.io/atlas` の直接依存が不要になれば除去する（`atlas-provider-gorm` は残る）。
- `Dockerfile` の `migrate` ビルドと `migrations/` の COPY を削除する（Cloud Run Job で適用しなくなるため）。
- `mise.toml` の `migrate:apply` タスクは `atlas migrate apply --env <env>` に置き換える。`migrate:diff` は据え置く。
- 権限付与（`apply-table-privileges`）の削除は、上記「3. 権限付与」の整理後に行う。

### 5. 移行手順

1. 既存環境の `atlas_schema_revisions` の中身が Atlas CLI から読めることを確認する（自作実装はテーブル定義を互換にしている）。差異があれば `atlas migrate set` で合わせる。
2. dry-run のみのワークフローを先に導入し、全環境で差分が空・想定どおりであることを確認する。
3. dev で適用を CI に切り替え、続いて stg、最後に prd へ展開する。
4. 安定後に `cmd/migrate` 等を削除する。

### 6. 未確定事項

- 現在、本番で `migrate` / `apply-table-privileges` Job を誰がいつ実行しているか（Terraform、手動 `gcloud run jobs execute` など）の確認が必要。
- GitHub Actions の WIF サービスアカウントに、各環境の Cloud SQL への接続権限（`roles/cloudsql.client` + DB IAM ユーザー）を付与する必要がある。権限は dry-run 用に読み取り専用のロールで分けるのが望ましい。
- Cloud SQL が Private IP のみの場合、GitHub-hosted runner から接続できない。その場合は Public IP + Auth Proxy、または Cloud Run Job での dry-run 実行に切り替える（`atlas` 公式イメージを使い、結果をログで確認する）。
- ファイル名が不揃いなマイグレーション（名前なし）の扱い。履歴の書き換えは不要で、今後は `mise run migrate:diff <name>` で名前を必須にする運用で足りる。
