# migrations

Atlas + versioned SQL でデータベーススキーマを管理する。GORM の AutoMigrate は完全に廃止し、`internal/shared/model/` の GORM モデルを「望ましいスキーマの宣言ソース」として残しつつ、Atlas の GORM provider で desired state を生成し、差分から SQL ファイルを起こす運用に統一する。

## 構成

- `internal/shared/model/` — 望ましいスキーマ (GORM モデル群)
- `atlas.hcl` — Atlas プロジェクト設定 (env "local" / env "prod"、dev DB は使い捨て docker postgres)
- `migrations/` — 生成された versioned SQL と `atlas.sum` を格納するディレクトリ
- `cmd/migrate-job/` — Cloud Run Job として Atlas を実行するエントリ。cloudsqlconn による IAM 認証付き `*sql.DB` 経由で `migrations/` を適用する

## ツール

開発者ローカル / CI で必要なものはすべて `mise.toml` にピンしている。`mise install` で揃う。

- `atlas` (CLI)
- `go` (atlas-provider-gorm を `go run` で呼ぶため)
- `task`

## 開発フロー

### 1. モデルを編集する

`internal/shared/model/<table>.go` を編集または新規追加する。

### 2. 差分を生成する

```bash
task migrate:diff -- <name>
```

内部的には次が走る。

```bash
atlas migrate diff <name> --env local
```

- 使い捨て docker postgres (`docker://postgres/16/dev`) が dev DB として起動する
- `internal/shared/model/` を atlas-provider-gorm が読み込み、desired state を SQL に落とす
- 既存の `migrations/<前回まで>.sql` を順に dev DB へ適用し、desired state との差分を `migrations/<timestamp>_<name>.sql` として書き出し、`atlas.sum` を更新する

### 3. 差分を lint する

```bash
task migrate:lint
```

PR では `migrate-lint.yml` ワークフローが同等の `atlas migrate lint --env local --latest 1` を必須実行する (詳細は `.github/workflows/migrate-lint.yml`)。

### 4. ローカル PostgreSQL で適用する (任意)

```bash
DEV_DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/dotto?sslmode=disable \
    task migrate:apply:local
```

### 5. PR レビュー

- `migrations/<timestamp>_*.sql` と `atlas.sum` の両方を必ず diff レビュー対象に含める
- 破壊的変更 (DROP / ALTER COLUMN TYPE 等) は `atlas migrate lint` の警告に従い、追加 migration や手動 hint で吸収する

## 本番適用

GitHub Actions `deploy.yml` が次の順で動く。

1. `terraform apply` で Cloud Run Service / Job (migrate-job 含む) のイメージ tag を更新
2. `gcloud run jobs execute migrate-job<env-suffix> --wait` を実行
3. Cloud Run Job 内で `cmd/migrate-job` が起動し、`cloudsqlconn` 経由で IAM 認証付き接続を確立後、`migrations/` を適用する

`migrate-job` のタイムアウトは `MIGRATE_TIMEOUT` (デフォルト `10m`) で外側から強制終了される。失敗時は non-zero exit でワークフローも失敗扱い。

## baseline (Phase 1 cutover)

Phase 1 の cutover では `internal/shared/model/` の GORM モデルを **正** として、初期 baseline SQL を起こす。手順は次の通り。

```bash
# 1. baseline 候補を生成
task migrate:diff -- baseline   # → migrations/<ts>_baseline.sql + atlas.sum

# 2. grants を補完
#    UUID 採番は Postgres 13+ 組み込みの gen_random_uuid() を使うため
#    uuid-ossp 拡張は不要。GRANT 系は ALTER DEFAULT PRIVILEGES FOR ROLE
#    dotto_admin GRANT ... に集約し、今後のテーブル追加でも自動的に GRANT が
#    伝播する形にする。
#    例:
#       ALTER DEFAULT PRIVILEGES FOR ROLE dotto_admin
#           IN SCHEMA public
#           GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO dotto_service;
#       ALTER DEFAULT PRIVILEGES FOR ROLE dotto_admin
#           IN SCHEMA public
#           GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO dotto_developer;

# 3. 本番 / stg / dev DB との突き合わせ
atlas schema diff \
    --from "postgres://<dotto_admin>@<host>/<db>?sslmode=disable" \
    --to file://migrations \
    --dev-url "docker://postgres/16/dev"

# 4. 差分があれば baseline を補正、または追加 migration で吸収する

# 5. 本番 / stg / dev DB に baseline を「適用済み」としてマーク
atlas migrate set --env prod \
    --url "<本番接続 URL>" \
    <baseline-version>
```

> **注意**: Cloud SQL は IAM 認証が前提のため、ローカルから atlas CLI で接続する場合は Cloud SQL Auth Proxy 経由で `127.0.0.1` を経由する必要がある。サービス本番では `cmd/migrate-job` が `cloudsqlconn` 経由で接続するためサイドカーは不要。

## ロールバック方針

- migration は **常に前進のみ**。明示的な DOWN を持たず、必要があれば「打ち消し migration」を新たに発行する
- 本番で問題が起きた場合は `terraform apply` で前回コミットの image tag に巻き戻し、追加 migration を適用するまで該当 PR を止める
- Cloud Run Job の `MIGRATE_TIMEOUT` をハングアップ防止用に必ず設定しておくこと
