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
- **GitHub Actions から Cloud SQL に IAM 認証で接続する。** Workload Identity Federation で認証し、Cloud SQL Auth Proxy 経由で接続する。環境ごとの GitHub Environment の `vars` で接続先（DB 名）を切り替える。インスタンスは全環境共通で、環境差は DB 名のみ。
- **PR で静的検証と dry-run を行う。** DB を変更しない検証をマージ前に必ず通す。

### ワークフロー

| トリガー | 環境 | 実行内容 |
| --- | --- | --- |
| `pull_request` / `merge_group` | dev, stg, prd（matrix） | `atlas migrate validate` → `atlas migrate lint` → モデル差分チェック → `atlas migrate apply --dry-run` |
| `push`（main） | dev, stg（matrix） | `atlas migrate apply` |
| `workflow_dispatch` | 入力で選択（prd 用） | `atlas migrate apply`（GitHub Environment `prd` の承認を必須にする） |

補足:

- dry-run は DB に書き込まない。ただし接続して未適用のリビジョンを読むため、PR の CI にも各環境への読み取り接続が必要になる。**PR の dry-run ジョブは GitHub Environment（特に `prd`）を指定しない。** `prd` は承認者必須かつ保護ブランチのみ許可（`protected_branches_only`）のため、指定すると PR のたびに承認待ちになる、または PR ブランチから実行できない。環境の切り替えは matrix の値（DB 名）で行う。fork からの PR ではシークレットが渡らないため、dry-run は同一リポジトリの PR のみ対象とする。
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

# stg / prd も同様（DB 名は GitHub Environment ごとの vars / PR 時は matrix の値で切り替える）
```

`DATABASE_URL` は `postgres://<iam-user>@127.0.0.1:5432/<dotto | dotto_stg | dotto_dev>?sslmode=disable`。インスタンスは共通で、環境ごとに変わるのは DB 名（と IAM ユーザー）のみ。

`env "gorm"`（diff 生成用）は据え置く。lint 用に `lint { latest = 1 }` などの設定も追加する。

### 2. ワークフロー

- 新規 `.github/workflows/migrate.yml` を追加する。`deploy.yml` とは別に管理し、トリガーと環境の組み合わせを上表のとおり定義する。
- 共通処理は composite action（例: `.github/actions/atlas-migrate`）に切り出す。内容は次のとおり。
  1. WIF で GCP 認証
  2. Cloud SQL Auth Proxy 起動（IAM 認証）
  3. mise で Atlas をセットアップ（`mise.toml` のバージョンと同一）
  4. `atlas migrate apply --env <env>`（dry-run 時は `--dry-run`）
- PR 向けの静的検証ジョブでは、dev DB として `docker://postgres/18/dev` を使う（`atlas.hcl` の既存設定）。
- 新しいチェック名（検証ジョブ・dry-run ジョブ）は infra 側の Ruleset の `required_status_checks`（現状 `Test` / `Lint` / `Build`）に追加する必要がある（「8. infra 側で必要な対応」参照）。matrix ジョブのチェック名は `ジョブ名 (env)` になる点に注意する。
- `concurrency` を環境ごとに設定し、同一環境へのマイグレーションの並列実行を防ぐ（Atlas は advisory lock も取得する）。
- `deploy.yml` の push 対象を stg から dev + stg に拡張するかは別途判断する。マイグレーションのみ dev / stg 両方を対象とし、アプリのデプロイ範囲は現状維持でもよい。

### 3. 権限付与

`apply-table-privileges`（`internal/migrate/privileges.go`）は、`current_schema()` 配下の全テーブルについて次を毎回やり直している。

- owner を `dotto_admin` に変更し、`ALL PRIVILEGES` を付与
- `dotto_service` / `dotto_developer` に `SELECT, INSERT, UPDATE, DELETE` を付与

今回の移行で、これを別ジョブとして残さず **マイグレーション適用に内包する**。方針は次のとおり。

1. **新規テーブルの owner を `dotto_admin` にする。** 適用用 SA（Cloud SQL IAM ユーザー）を `dotto_admin` ロールのメンバーにし、接続時に `role=dotto_admin`（`SET ROLE` 相当）を指定する。`DATABASE_URL` に `?options=-c%20role%3Ddotto_admin` を付ける案、または各マイグレーションの先頭に `SET ROLE dotto_admin;` を書く案がある。前者はファイルを汚さないが、Atlas（pgx）が実行時パラメータとして透過することの検証が必要。
2. **新規テーブルへの権限は `ALTER DEFAULT PRIVILEGES` で自動付与する。** `dotto_admin` が作るオブジェクトに対し `dotto_service` / `dotto_developer` へ上記権限が付くよう、DB ごとに 1 回設定する。
3. **既存テーブルの owner 変更・権限付与は、マイグレーションに含めず手動で行う（担当: kantacky）。** 各 DB（`dotto` / `dotto_stg` / `dotto_dev`）で、現行の `apply-table-privileges` と同じ内容（owner を `dotto_admin` に変更、`dotto_admin` に `ALL`、`dotto_service` / `dotto_developer` に `SELECT, INSERT, UPDATE, DELETE`）を実施する。あわせて 2 の `ALTER DEFAULT PRIVILEGES` も同じタイミングで 1 回実行する（マイグレーションには入れない）。

```sql
-- 手動実行用（各 DB で 1 回、dotto_admin 権限で実行）
ALTER DEFAULT PRIVILEGES FOR ROLE dotto_admin
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO dotto_service, dotto_developer;
```

手動対応が完了するまでは、適用ワークフロー（`push` / `workflow_dispatch`）を有効化しない。

注意点:

- `ALTER DEFAULT PRIVILEGES` は **DB 単位**（`dotto` / `dotto_stg` / `dotto_dev` それぞれ）かつ **作成者ロール単位**。`dotto_admin` 以外のロールで作ったテーブルには効かないため、1 の `role=dotto_admin` が前提になる。
- 現行ジョブはシーケンスへの権限を扱っていない。シーケンスを使うテーブル（`serial` / `bigserial`）があれば `USAGE` の付与も必要になる。要確認。
- 人手で作ったテーブルやマイグレーション外の変更は、この方式では権限が補正されない（現行ジョブの「全テーブル再付与」と挙動が異なる）。必要なら移行後も `mise` タスクとして手動で流せるようにしておく。

これにより `apply-table-privileges` と Cloud Run Job は不要になり、「4. コードの削除」で一緒に削除できる。

### 4. コードの削除

- `cmd/migrate`、`internal/migrate/{apply,revisions}.go` を削除する。
- `go.mod` から `ariga.io/atlas` の直接依存が不要になれば除去する（`atlas-provider-gorm` は残る）。
- `Dockerfile` の `migrate` / `apply-table-privileges` ビルドと `migrations/` の COPY を削除する（Cloud Run Job で適用しなくなるため）。
- `mise.toml` の `migrate:apply` タスクは `atlas migrate apply --env <env>` に置き換える。`migrate:diff` は据え置く。
- `cmd/apply-table-privileges`、`internal/migrate/privileges.go`、`Dockerfile` の `apply-table-privileges` ビルド、`set-env` のコメント中の記載を削除する。「3. 権限付与」の方式で全環境の権限が揃ったことを確認してから行う。

### 5. 移行手順

1. 既存環境の `atlas_schema_revisions` の中身が Atlas CLI から読めることを確認する（自作実装はテーブル定義を互換にしている）。差異があれば `atlas migrate set` で合わせる。
2. dry-run のみのワークフローを先に導入し、全環境で差分が空・想定どおりであることを確認する。
3. 既存テーブルの権限変更と `ALTER DEFAULT PRIVILEGES` を各 DB で手動実施する（dev → stg → prd）。実施後、現行ジョブの結果と owner / 権限が一致することを確認する。
4. dev で適用を CI に切り替え、続いて stg、最後に prd へ展開する。
5. 安定後に `cmd/migrate` 等を削除する。

### 6. Cloud SQL への IAM 認証接続

Cloud SQL は IAM 認証で接続している（現状の `cmd/migrate` も cloudsqlconn + IAM 認証）。Atlas CLI からは Cloud SQL Auth Proxy の `--auto-iam-authn` を介して、パスワードなしで接続する。

```bash
cloud-sql-proxy --auto-iam-authn <PROJECT>:<REGION>:<INSTANCE> &
export DATABASE_URL="postgres://<iam-user>@127.0.0.1:5432/<db>?sslmode=disable"
atlas migrate apply --env <env> [--dry-run]
```

- プロキシが OAuth トークンをパスワード代わりに渡す。プロキシ〜インスタンス間は暗号化されるため、ローカル側は `sslmode=disable` でよい。
- GitHub Actions では `google-github-actions/auth`（WIF）→ プロキシをバックグラウンド起動 → Atlas 実行の順に composite action へまとめる。
- IAM アクセストークンの有効期間は約 1 時間。マイグレーションは短時間で終わるため問題ない。

#### 必要な設定

| 項目 | 内容 |
| --- | --- |
| DB ユーザー | WIF 用サービスアカウントを Cloud SQL の IAM ユーザーとして登録する（ユーザー名は `<sa-name>@swift2023groupc.iam`、`.gserviceaccount.com` を除いた形）。infra では `terraform/google/foundation/cloud-sql` の `google_sql_user` で管理しており、`database_service_accounts` に追加して `database_roles` を指定する |
| IAM ロール | `roles/cloudsql.client`、`roles/cloudsql.instanceUser` |
| DB 権限 | dry-run 用: 接続と `atlas_schema_revisions` の読み取り。適用用: `dotto_admin` ロールのメンバーとし、`role=dotto_admin` で接続して DDL を実行する（「3. 権限付与」） |

現状の DB ロール（infra 確認済み）: アプリ・ジョブの SA は `dotto_service`、人間の管理者は `dotto_admin`、開発者は `dotto_developer`。マイグレーション用 SA は `dotto_service` のままでは DDL ができないため、適用用は `dotto_admin` ロールを割り当てる。

#### サービスアカウントの分離

- **dry-run 用（読み取り専用）:** `pull_request` / `merge_group` で使う。fork からの PR には渡さない。
- **適用用（DDL 権限あり）:** `push`（main）と `workflow_dispatch` で使う。GitHub Environment ごとに別の SA / 変数を割り当て、prd は承認必須にする。

### 7. 未確定事項

- 現在、本番で `migrate` / `apply-table-privileges` Job を誰がいつ実行しているか。**Terraform に Job 定義がない**ことは確認済みなので、手動作成・手動実行（`gcloud run jobs execute` 等）の可能性が高い。GCP 上の実在確認が必要。
- GitHub Actions 用の接続権限の付与方式（「8. infra 側で必要な対応」で整理）。
- `dotto_dev` / `dotto_stg` / `dotto` の各 DB のスキーマ・`atlas_schema_revisions` の状態（infra からは確認できない）。
- ファイル名が不揃いなマイグレーション（名前なし）の扱い。履歴の書き換えは不要で、今後は `mise run migrate:diff <name>` で名前を必須にする運用で足りる。

### 8. infra 側で必要な対応

`fun-dotto/infra` を確認した結果、以下が必要。いずれも Terraform（Digger 経由）で管理する。

#### 8.1 現状（確認済み）

- `terraform/github/fun-dotto/repository-server.tf`: `server` リポジトリの `workload_identity_roles` は `artifactregistry.writer` / `iam.serviceAccountUser` / `run.developer` のみ。**Cloud SQL 関連のロールはない。** GitHub Environment は `dev` / `stg` / `prd` が定義済みで、`prd` のみ承認者 3 名・`prevent_self_review`・`protected_branches_only`・`can_admins_bypass = false`。
- WIF は `terraform/modules/github-repository/workload-identity.tf` で **リポジトリ単位の principal**（`principal://.../subject/<owner>/<repo>`）に直接ロールを付与する方式。プロバイダの属性マッピングは `google.subject = assertion.repository` のみで、**環境（`assertion.environment`）では区別できない。** つまり現構成のままでは、`server` リポジトリのどのワークフローも同じ権限を持つ。
- `terraform/google/foundation/cloud-sql`: IAM ユーザーは `database_service_accounts`（アプリ・ジョブ用 SA、`dotto_service`）、`database_admins`（`dotto_admin`）、`database_developers`（`dotto_developer`）。**マイグレーション用の SA / DB ユーザーはない。**
- `terraform/google/foundation/iam-workload-identity` に `terraform-ci` / `github-actions-app` の SA がある。リポジトリ専用 SA を `workloadIdentityUser` で借用するパターンが既にある（直近の PR で app リポジトリ専用 SA を追加済み）。
- `required_status_checks` は `["Test", "Lint", "Build"]`。

#### 8.2 対応一覧

| # | 対応 | 場所 | 備考 |
| --- | --- | --- | --- |
| 1 | マイグレーション用 SA を作成する（dry-run 用 / 適用用の 2 種） | `terraform/google/foundation/iam-workload-identity` | 既存の `github-actions-app` のパターンに倣う。`server` リポジトリの principal に `roles/iam.workloadIdentityUser` を付与 |
| 2 | 各 SA に `roles/cloudsql.client` / `roles/cloudsql.instanceUser` を付与 | 同上 | プロジェクトレベル。Cloud SQL は 1 インスタンスのため環境別の IAM 条件は不要だが、DB 単位の権限は 3 に委ねる |
| 3 | 各 SA を Cloud SQL の IAM ユーザーとして登録し、DB ロールを割り当てる | `terraform/google/foundation/cloud-sql`（`locals.tf` / `users.tf`） | 現状の `database_service_accounts` は `dotto_service` 固定。適用用は `database_roles = ["dotto_admin"]`（「3. 権限付与」の `role=dotto_admin` 接続に必要）、dry-run 用は読み取り専用ロール（新規 DB ロール）が必要。`google_sql_user` は `database_roles` を指定できるが、**ロール自体の GRANT（`atlas_schema_revisions` の SELECT、`dotto_dev` / `dotto_stg` / `dotto` ごとの権限）は Terraform の管理外の SQL** になる可能性が高く、手順化が必要 |
| 4 | 環境ごとの権限分離をどう担保するか決める | WIF / SA 設計 | リポジトリ単位 principal では環境を区別できないため、(a) プロバイダに `attribute.environment` を追加して SA 借用を環境条件で制限する、(b) 環境ごとに SA を分け Environment の `vars` で SA を切り替える（ただし PR からも他環境の SA を指定できてしまう）、のいずれか。prd 適用 SA は `environment:prd` かつ保護ブランチのみ借用可とするのが望ましい。(a) は全リポジトリの WIF に影響するため要検討 |
| 5 | `server` リポジトリの GitHub Environment に変数を追加する | `terraform/github/fun-dotto`（`github_actions_environment_variable` を module に追加） | DB 名（`dotto` / `dotto_stg` / `dotto_dev`）、インスタンス接続名、使用する SA、WIF プロバイダ名など。現在 module に Environment 変数のリソースはない |
| 6 | Ruleset の `required_status_checks` に新しいチェックを追加する | `repository-server.tf` | 検証ジョブ・dry-run ジョブ（matrix なら `名前 (env)` 形式）。`merge_group` でも同名チェックが走ること。ワークフロー導入後に追加しないとマージ不能になるため、順序に注意 |
| 7 | `prd` Environment の見直し | `repository-server.tf` | 既存設定（承認者必須・`protected_branches_only`）は適用ワークフローの要件に合致。`dev` / `stg` は現状保護なしで、`push`（main）での適用に問題はない。PR の dry-run は Environment を使わない（本文参照）ため追加変更は不要 |
| 8 | 既存の `migrate` / `apply-table-privileges` Job の実在確認と削除 | GCP コンソール / `gcloud run jobs list` | Terraform 管理外。権限付与をマイグレーションに内包するため、移行完了後に両方削除する（ワークフローから起動しないので `run.developer` の追加は不要） |
| 8b | `dotto_dev` / `dotto_stg` / `dotto` の各 DB で、適用用 SA が `dotto_admin` に `SET ROLE` できることを確認する | Cloud SQL（SQL 手動） | `google_sql_user` の `database_roles` でメンバーになるが、`dotto_admin` が `INHERIT` / `SET` オプション付きのメンバーシップを許すか（PostgreSQL 16+ の `GRANT ... WITH SET TRUE`）を実機で検証する |
| 9 | ワークフロー導入後の drift 確認 | `.github/workflows/drift-detection.yml` | 追加リソースが Terraform 管理に入っていることの確認のみ |

#### 8.3 推奨する進め方

1. infra: 8.2 の 1〜3, 5（SA・IAM ユーザー・Environment 変数）を dev から先行して入れる。DB ロールの GRANT は手動 SQL になるため、実施手順を infra リポジトリのドキュメントに残す。
2. server: dry-run のみの `migrate.yml` を追加し、動作確認後に infra で 8.2 の 6（必須チェック追加）を行う。
3. 適用ワークフローを dev → stg → prd の順に有効化する。prd の SA の借用条件（8.2 の 4）は prd 展開前に必ず確定させる。
