# infra

`fun-dotto/server` モノレポの Terraform 構成。Cloud Run Service / Job、Artifact Registry、Service Account、Cloud Scheduler を `for_each` で集約し、環境は `envs/<env>.tfvars` で `-var-file` 渡しする。

- backend: `gcs` バケット `swift2023groupc-tfstate`、prefix `server`
- provider: `hashicorp/google ~> 6.0`
- Terraform: `1.9.8` (mise.toml でピン)
- 命名規約: prod は `<resource>`、それ以外は `<resource>-<env>` (way の set-env と同形)

## ディレクトリ

| ファイル | 役割 |
| --- | --- |
| `main.tf` | backend 設定 / provider 宣言 / 必要な Google API 有効化 |
| `variables.tf` | 環境変数定義 (project_id / region / environment / DB / image_tag / cron) |
| `locals.tf` | env_suffix と http_services / cloud_run_jobs マップ |
| `artifact_registry.tf` | `server` Docker repo (1 イメージに全バイナリを同梱) |
| `service_account.tf` | Service / Job ごとの SA、cloudsql.client/instanceUser、google_sql_user、fcm_sender カスタムロール、scheduler SA |
| `cloud_run_services.tf` | HTTP サービス (`for_each = local.http_services`) |
| `cloud_run_jobs.tf` | Cloud Run Job (`for_each = local.cloud_run_jobs`、migrate-job も含む) |
| `scheduler.tf` | Cloud Scheduler (`for_each = local.scheduled_jobs`) と Job invoker IAM |
| `iam.tf` | HTTP サービスの公開 (allUsers invoker) |
| `outputs.tf` | image / service URL / job name / SA email マップ |
| `envs/<env>.tfvars` | 環境ごとの変数値 |

## 初回ブートストラップ

Cloud Run Service / Job を立てる前に Artifact Registry と SA / Cloud SQL ユーザを先に作る必要がある (chicken-and-egg 回避)。`batch-jobs/infra` と同じ手順を踏む。

```bash
# 1. プロジェクト・GCS state バケット・Cloud SQL インスタンスは事前準備済み (本リポジトリの管理対象外)
# 2. Workload Identity Federation の設定 (way 配下のドキュメント参照)

# 3. AR / SA / IAM だけ先に apply
terraform init
terraform apply \
    -var-file=envs/dev.tfvars \
    -target=google_project_service.required_apis \
    -target=google_artifact_registry_repository.server \
    -target=google_service_account.workload \
    -target=google_project_iam_member.workload_sql_client \
    -target=google_project_iam_member.workload_sql_instance_user \
    -target=google_sql_user.workload \
    -target=google_project_iam_custom_role.fcm_sender \
    -target=google_project_iam_member.dispatch_notifications_fcm_sender \
    -target=google_service_account.scheduler \
    -target=google_service_account_iam_member.scheduler_token_creator

# 4. 1 度 Docker イメージを push (GitHub Actions deploy.yml で自動化される)

# 5. 残りを apply
terraform apply -var-file=envs/dev.tfvars -var=image_tag=<sha>
```

## 既存 academic-api Cloud Run Service の取り込み (cutover)

cutover 段階で `academic-api` / `academic-api-stg` / `academic-api-dev` Cloud Run Service を server state に取り込む。リソース名・リージョン・SA・環境変数キーを既存に合わせて差分ゼロにする。

```bash
terraform import \
    -var-file=envs/dev.tfvars \
    'google_cloud_run_v2_service.http["academic-api"]' \
    "projects/<project_id>/locations/asia-northeast1/services/academic-api-dev"
```

prod / stg も同様に実施し、`terraform plan` が clean になるまで属性を合わせ込む (詳細は Notion 計画 §H)。

## Terraform 対象外として明記

- Cloud SQL インスタンス本体: 既存を `instance_connection_name` 経由で参照のみ
- `dotto_admin` / `dotto_service` / `dotto_developer` ロールおよび `uuid-ossp` 拡張: 手動セットアップ前提 (`shared-go/db/migrate.go` 由来の規約)
- 旧 Artifact Registry repo (academic-api / batch-jobs Go ジョブ用): Cleanup Policy に任せる
