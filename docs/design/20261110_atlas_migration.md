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

実装では、dev / stg の `push`（main）の適用は、アプリのデプロイより前に完了させるため `deploy.yml` から呼び出す（「ワークフロー」参照）。

## 移行前の現状

### 構成

| 役割 | 実装 |
| --- | --- |
| スキーマの源泉 | GORM モデル（`internal/shared/model`）を `cmd/atlas-loader` が列挙し、`atlas-provider-gorm` で SQL 化 |
| diff 生成 | `mise run migrate:diff`（`atlas migrate diff --env gorm`）。dev DB は `docker://postgres/18/dev` |
| 適用 | `mise run migrate:apply`（`go run ./cmd/migrate`）。Atlas CLI は使わず、Atlas の Go API（`migrate.NewExecutor`）で適用 |
| 本番適用 | 同一イメージの `/bin/migrate` を Cloud Run Job で実行。cloudsqlconn + IAM 認証で接続 |
| リビジョン管理 | `internal/migrate/revisions.go` に Atlas CLI 互換の `atlas_schema_revisions` を自前実装（約 200 行） |
| 権限 | 別ジョブ `apply-table-privileges` が全テーブルの owner を `dotto_admin` に変更し、権限を再付与 |
| DB 構成（infra 確認済み） | Cloud SQL インスタンスは 1 つ（`swift2023groupc:asia-northeast1:dotto`、PostgreSQL 18、**Public IP 有効**、`cloudsql.iam_authentication=on`）。環境の分離は **同一インスタンス内の DB 名**（prd: `dotto` / stg: `dotto_stg` / dev: `dotto_dev`） |
| `atlas.hcl` | `env "gorm"` のみ（`src` + `dev` + `migration.dir`）。接続先 `url` を持つ env はない |

### CI / CD

- `ci.yml`: test / build / lint のみ。マイグレーションの検証はない。
- `deploy.yml`: main への push で stg にデプロイ。`workflow_dispatch` で環境を選択でき、prd は `action: deploy` のみ指定可。**`migrate` / `apply-table-privileges` の Job を実行するステップはない**（`set-env` の `DEPLOY_TARGETS` にも含まれない）。適用は手動、または別管理と思われる。
- infra 側（確認済み）: `migrate` / `apply-table-privileges` を実行する Cloud Run Job は **Terraform に定義がない**（`terraform/modules/cloud-run-job` を使う Job は import-gtfs / notify-user / build-irregular-notification のスケジュール実行系のみ）。手動作成または未管理の可能性が高い。
- dev へは `workflow_dispatch` でのみデプロイ可能。push 時の対象は stg のみ。

### 課題

1. **dry-run できない。** `cmd/migrate` は常に適用する。`atlas.hcl` に接続先を持つ env がないため `atlas migrate apply --dry-run` も使えない。
2. **PR で検証されない。** `atlas migrate lint` / `validate`（`atlas.sum` 整合性）や、モデルとマイグレーションの差分残りの検査がない。生成し忘れや破壊的変更に PR で気付けない。
3. **適用が CI/CD に組み込まれていない。** 環境ごとの適用トリガーが定義されておらず、実行順序（マイグレーション → アプリのデプロイ。権限付与はマイグレーションに内包する）の保証もない。
4. **Atlas CLI の再実装を抱えている。** リビジョンテーブルの読み書きを自前実装しており、Atlas のバージョンアップで互換性が崩れるリスクがある。ロックもなく、Job の多重起動に弱い。
5. Atlas の CLI（mise: 1.3.0）と Go ライブラリ（`ariga.io/atlas`）の二重管理。
6. マイグレーションファイル名が不揃い（`20260510085829.sql` のように名前なしのものがある）。

## 理想形

### 方針

- **適用も Atlas CLI に一本化する。** `atlas migrate apply` を使い、`--dry-run` で適用予定を確認できるようにする。自作のリビジョン管理（`revisions.go`）と `cmd/migrate` は廃止する。
- **環境ごとに `atlas.hcl` の env を定義する。** 接続先は環境変数（`DATABASE_URL`）で注入する。
- **GitHub Actions から Cloud SQL に IAM 認証で接続する。** Workload Identity Federation でサービスアカウントを借用し、Cloud SQL Auth Proxy 経由で接続する。インスタンスは全環境共通で、環境差は DB 名のみ。
- **PR で静的検証と dry-run を行う。** DB を変更しない検証をマージ前に必ず通す。dry-run の結果は PR にコメントする。

### ワークフロー

| トリガー | ワークフロー | 環境 | 実行内容 |
| --- | --- | --- | --- |
| `pull_request` / `merge_group` | `migrate.yml` の `Migration Validate` | 環境に依存しない（1 回） | `atlas migrate validate` → モデル差分チェック |
| `pull_request` / `merge_group` | `migrate.yml` の `Migration Dry Run` | dev, stg, prd（matrix） | `atlas migrate apply --dry-run`（PR にコメント） |
| `push`（main） | `deploy.yml` の `migrate` → `migrate-apply.yml` | dev, stg | `atlas migrate apply`。完了後にアプリをデプロイ（`needs`） |
| `workflow_dispatch`（Deploy） | `deploy.yml` の `migrate` → `migrate-apply.yml` | 選択した環境（`build-push` と prd を除く） | `atlas migrate apply`。完了後にアプリをデプロイ |
| `workflow_dispatch`（Migrate） | `migrate.yml` の `apply` → `migrate-apply.yml` | 入力で選択（prd 用） | `atlas migrate apply`（GitHub Environment `prd` の承認が必須） |

補足:

- dry-run は DB に書き込まない。ただし接続して未適用のリビジョンを読むため、PR の CI にも各環境への読み取り接続が必要になる。**PR の dry-run ジョブは GitHub Environment（特に `prd`）を指定しない。** `prd` は承認者必須かつ保護ブランチのみ許可（`protected_branches_only`）のため、指定すると PR のたびに承認待ちになる、または PR ブランチから実行できない。環境の切り替えは matrix の値（DB 名）で行う。fork からの PR ではシークレットが渡らないため、dry-run は同一リポジトリの PR のみ対象とする。
- 静的検証（validate / 差分チェック）は環境に依存しないため 1 回だけ実行し、環境ごとの matrix は `apply --dry-run` のみに絞る。
- `merge_group` でも同じ検証を実行し、マージキュー上の最終状態を検証する。
- `atlas migrate lint` は Atlas v0.38 以降ログイン（有償機能の認証）が必要なため採用しない。破壊的変更はレビューで確認する。
- **prd のマイグレーションは `deploy.yml` からは適用しない。** `deploy.yml` の prd は検証済みのイメージタグを指定する `deploy` のみで、デプロイするイメージとリポジトリの HEAD のマイグレーションがずれる可能性があるため。prd は `Migrate` の `workflow_dispatch` で、dev / stg で適用済みのコミットに対して、デプロイより前に手動で適用する。
- アプリのデプロイより前にマイグレーションが完了するよう、`deploy.yml` の `deploy` ジョブは `migrate` ジョブに `needs` する。マイグレーションが失敗した場合はデプロイしない。

## 実装

### 1. `atlas.hcl`

環境ごとの env を定義する。接続先は `getenv` で注入する。

```hcl
env "dev" {
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}

# stg / prd も同様
```

`DATABASE_URL` は `postgres://<iam-user>@127.0.0.1:5432/<dotto | dotto_stg | dotto_dev>?sslmode=disable`（IAM ユーザー名の `@` は `%40` にエスケープする）。適用時は `&options=-c%20role%3Ddotto_admin` を付け、接続後に `dotto_admin` へ切り替える。インスタンスは共通で、環境ごとに変わるのは DB 名のみ。

`env "gorm"`（diff 生成用）は据え置く。

### 2. ワークフロー

- `.github/workflows/migrate.yml`: 検証（`Migration Validate`）、dry-run（`Migration Dry Run`）、手動適用（`apply`）。
- `.github/workflows/migrate-apply.yml`: 適用の再利用ワークフロー（`workflow_call`）。入力の環境の配列を、環境ごとに GitHub Environment を指定して 1 つずつ適用する。`deploy.yml` と `migrate.yml` から呼び出す。
- `.github/actions/atlas-migrate`: 共通処理の composite action。
  1. WIF でサービスアカウントを借用して GCP 認証
  2. mise で Atlas をセットアップ（`mise.toml` のバージョンと同一）
  3. Cloud SQL Auth Proxy 起動（`--auto-iam-authn`）
  4. `atlas migrate apply --env <env>`（dry-run 時は `--dry-run`）。出力は PR コメント用にファイルにも保存する
- 新しいチェック名は infra 側の Ruleset の `required_status_checks` に登録済み（`Migration Validate`、`Migration Dry Run (dev, dotto_dev)`、`Migration Dry Run (stg, dotto_stg)`、`Migration Dry Run (prd, dotto)`）。matrix ジョブのチェック名は `ジョブ名 (environment, database_name)` になるため、matrix の値を変える場合は infra 側も更新する。
- `concurrency` を環境ごとに設定し、同一環境へのマイグレーションの並列実行を防ぐ（Atlas は advisory lock も取得する）。
- 適用ジョブを一時的に無効化するスイッチ（リポジトリ変数）は設けていない。権限の手動対応が終わるまでは、マージを控えることで制御する。

### 3. 権限付与

`apply-table-privileges` は、`current_schema()` 配下の全テーブルについて次を毎回やり直していた。

- owner を `dotto_admin` に変更し、`ALL PRIVILEGES` を付与
- `dotto_service` / `dotto_developer` に `SELECT, INSERT, UPDATE, DELETE` を付与

これを別ジョブとして残さず **マイグレーション適用に内包した**。方針は次のとおり。

1. **新規テーブルの owner を `dotto_admin` にする。** 適用用 SA（Cloud SQL IAM ユーザー）は Terraform の `database_roles` で `dotto_admin` のメンバーになっており（`inherit_option` / `set_option` がともに true）、接続時に `options=-c role=dotto_admin` を指定して切り替える。Atlas（pgx）がこのパラメータを透過することは、dev / stg への適用で確認済み。
2. **新規テーブルへの権限は `ALTER DEFAULT PRIVILEGES` で自動付与する。** `dotto_admin` が作るオブジェクトに対して付与される。
3. **既存テーブルの owner 変更・権限付与は、マイグレーションに含めず手動で行う。**

| ロール | dotto_dev / dotto_stg | dotto (prd) |
| --- | --- | --- |
| `dotto_admin` | owner、`ALL` | owner、`ALL` |
| `dotto_service` | `SELECT, INSERT, UPDATE, DELETE` | `SELECT, INSERT, UPDATE, DELETE` |
| `dotto_developer` | `SELECT, INSERT, UPDATE, DELETE` | **`SELECT` のみ** |

prd では運用上、`dotto_developer` には `SELECT` のみを与える。現行の `apply-table-privileges` は全環境で `dotto_developer` に DML を付与していたため、prd への切り替え時に既存テーブルの `INSERT, UPDATE, DELETE` を取り消す。手動実行する SQL は infra の `docs/atlas-migration.md` にまとめている。

注意点:

- `ALTER DEFAULT PRIVILEGES` は **DB 単位**（`dotto` / `dotto_stg` / `dotto_dev` それぞれ）かつ **作成者ロール単位**。`dotto_admin` 以外のロールで作ったテーブルには効かないため、1 の `role=dotto_admin` が前提になる。
- 現行ジョブはシーケンスへの権限を扱っていない。シーケンスを使うテーブル（`serial` / `bigserial`）が追加された場合は、`USAGE` の付与（`ALTER DEFAULT PRIVILEGES ... ON SEQUENCES`）も必要になる。
- 人手で作ったテーブルやマイグレーション外の変更は、この方式では権限が補正されない（旧ジョブの「全テーブル再付与」と挙動が異なる）。

stg では、GTFS のマイグレーションを適用して、新規テーブルの owner が `dotto_admin`、権限が上表どおりに付くことを確認した。

### 4. コードの削除

次を削除した。

- `cmd/migrate`、`cmd/apply-table-privileges`、`internal/migrate/`
- `go.mod` の `ariga.io/atlas` の直接依存（`atlas-provider-gorm` が使うため indirect として残る）
- `Dockerfile` の `migrate` / `apply-table-privileges` のビルドと `migrations/` の COPY
- `mise.toml` の `migrate:apply` を `atlas migrate apply --env <env>` に置き換え（`migrate:diff` は据え置き）

GCP 上の `migrate` / `apply-table-privileges` Cloud Run Job は Terraform の管理外のため、`gcloud run jobs list` で実在を確認して手動で削除する。

### 5. 移行手順

1. 既存環境の `atlas_schema_revisions` の中身が Atlas CLI から読めることを確認する。**完了**（stg / prd の dry-run で現在のバージョンを読み取れた）。
2. dry-run のみのワークフローを導入し、全環境で差分が想定どおりであることを確認する。**完了**。
3. 既存テーブルの権限変更と `ALTER DEFAULT PRIVILEGES` を各 DB で手動実施する（dev → stg → prd）。**dev / stg は完了。prd は適用前に実施する。**
4. dev で適用を CI に切り替え、続いて stg、最後に prd へ展開する。**dev / stg は完了。prd は未実施。**
5. 安定後に `cmd/migrate` 等を削除する。**完了**（上記）。

### 6. Cloud SQL への IAM 認証接続

Cloud SQL は IAM 認証で接続している。Atlas CLI からは Cloud SQL Auth Proxy の `--auto-iam-authn` を介して、パスワードなしで接続する。

```bash
cloud-sql-proxy --auto-iam-authn <PROJECT>:<REGION>:<INSTANCE> &
export DATABASE_URL="postgres://<iam-user>@127.0.0.1:5432/<db>?sslmode=disable"
atlas migrate apply --env <env> [--dry-run]
```

- プロキシが OAuth トークンをパスワード代わりに渡す。プロキシ〜インスタンス間は暗号化されるため、ローカル側は `sslmode=disable` でよい。
- IAM アクセストークンの有効期間は約 1 時間。マイグレーションは短時間で終わるため問題ない。

#### 必要な設定

| 項目 | 内容 |
| --- | --- |
| DB ユーザー | サービスアカウントを Cloud SQL の IAM ユーザーとして登録する（ユーザー名は `<sa-name>@swift2023groupc.iam`、`.gserviceaccount.com` を除いた形）。infra の `terraform/google/foundation/cloud-sql/users.tf` で管理 |
| IAM ロール | `roles/cloudsql.client`、`roles/cloudsql.instanceUser` |
| DB 権限 | dry-run 用: 接続と `atlas_schema_revisions` の読み取り。適用用: `dotto_admin` ロールのメンバーとし、`role=dotto_admin` で接続して DDL を実行する |

リビジョンテーブルは `public` ではなく、専用スキーマ `atlas_schema_revisions` の `atlas_schema_revisions` テーブルにある。dry-run 用ロールには、このスキーマの `USAGE` とテーブルの `SELECT` を付与する。

#### サービスアカウントの分離

- **dry-run 用（読み取り専用）`gha-server-migrate-dry-run`:** `pull_request` / `merge_group` で使う。server リポジトリの任意のワークフローから借用できる（読み取り専用のため）。fork からの PR には OIDC トークンが渡らない。
- **適用用（DDL 権限あり）`gha-server-migrate-apply`:** 適用ジョブで使う。server の GitHub Environment（dev / stg / prd）を指定したジョブからのみ借用できる。prd は Environment 側の承認者必須・保護ブランチ限定で守る。

### 7. infra 側の対応

`fun-dotto/infra` で次を実施した。いずれも Terraform（Digger 経由）で管理している。

| 対応 | 場所 | 状況 |
| --- | --- | --- |
| マイグレーション用 SA（dry-run 用 / 適用用）の作成 | `terraform/google/foundation/iam-workload-identity` | 完了（infra#100） |
| 各 SA への `roles/cloudsql.client` / `roles/cloudsql.instanceUser` の付与 | 同上 | 完了（infra#100） |
| SA の Cloud SQL IAM ユーザー登録とロール割り当て | `terraform/google/foundation/cloud-sql` | 完了（infra#100）。適用用は `dotto_admin` |
| 環境ごとの権限分離 | WIF プロバイダ | 完了（infra#100）。プロバイダに `attribute.repository_environment`（`<owner>/<repo>/<environment>`）を追加し、適用用 SA の借用を server の Environment 限定にした |
| `required_status_checks` の追加 | `terraform/github/fun-dotto/repository-server.tf` | 完了（infra#101） |
| DB ロールの GRANT（読み取り専用ロール、`ALTER DEFAULT PRIVILEGES` など） | Cloud SQL（手動 SQL） | dev / stg は完了、prd は適用前に実施する。手順は infra の `docs/atlas-migration.md` |
| GitHub Environment への変数追加 | — | 不要。DB 名などは PR の dry-run が Environment を使えないため、ワークフローの matrix とワークフロー内の `env` で持つ |
| 既存の `migrate` / `apply-table-privileges` Job の実在確認と削除 | GCP（`gcloud run jobs list`） | 未実施（手動） |

`required_status_checks` の追加や prd の `dotto_developer` の権限などの運用は、上記の infra の docs を参照する。

### 8. 残課題

- prd への適用（`Migrate` の `workflow_dispatch`）と、prd の権限の手動対応（`dotto_developer` を `SELECT` のみにする）。
- GCP 上の旧 Cloud Run Job（`migrate` / `apply-table-privileges`）の実在確認と削除。
- マイグレーションファイル名の不揃い（名前なしのものがある）。履歴の書き換えは不要で、今後は `mise run migrate:diff <name>` で名前を付ける運用にする。
