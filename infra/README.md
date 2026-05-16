# infra

`fun-dotto/server` モノレポの Terraform 構成。Cloud Run Service / Job、Artifact Registry、Service Account、Cloud Scheduler を `for_each` で集約し、環境は `envs/<env>.tfvars` で `-var-file` 渡しする。

- backend: `gcs` バケット `swift2023groupc-tfstate`、prefix は **env ごとに `server/<env>`** へ分離 (`terraform init -backend-config="prefix=server/<env>"`)
- provider: `hashicorp/google ~> 6.0`
- Terraform: `1.9.8` (mise.toml でピン)
- 命名規約: prod は `<resource>`、それ以外は `<resource>-<env>` (way の set-env と同形)

## ディレクトリ

| ファイル | 役割 |
| --- | --- |
| `main.tf` | backend 設定 / provider 宣言 / 必要な Google API 有効化 |
| `variables.tf` | 環境変数定義 (project_id / region / environment / DB / image_tag / cron) |
| `locals.tf` | env_suffix と http_services / cloud_run_jobs マップ、SA account_id 用 sa_id、fcm_sender role 固定パス |
| `artifact_registry.tf` | `server` Docker repo (1 イメージに全バイナリを同梱、prod state でのみ create) |
| `service_account.tf` | Service / Job ごとの SA、cloudsql.client/instanceUser、google_sql_user、fcm_sender カスタムロール (prod のみ)、scheduler SA |
| `cloud_run_services.tf` | HTTP サービス (`for_each = local.http_services`) |
| `cloud_run_jobs.tf` | Cloud Run Job (`for_each = local.cloud_run_jobs`、migrate-job も含む) |
| `scheduler.tf` | Cloud Scheduler (`for_each = local.scheduled_jobs`) と Job invoker IAM |
| `iam.tf` | HTTP サービスの公開 (allUsers invoker) |
| `outputs.tf` | image / service URL / job name / SA email マップ |
| `envs/<env>.tfvars` | 環境ごとの変数値 (project_id / instance_connection_name は `REPLACE_ME` プレースホルダ、実値は手元で上書き) |

## state 分離

backend prefix を `server` 固定にすると 4 env で 1 state を共有してしまい、`-var-file` を切り替えた瞬間に他環境のリソースを置換/削除する事故が起きる。env ごとに prefix を分けて apply する:

```bash
# dev
terraform init -reconfigure -backend-config="prefix=server/dev"
terraform apply -var-file=envs/dev.tfvars

# prod
terraform init -reconfigure -backend-config="prefix=server/prod"
terraform apply -var-file=envs/prod.tfvars
```

(`terraform workspace` で代替してもよいが、CI/CD では `-backend-config` で env を 1 ジョブ 1 state に固定する方が事故が少ない。)

## 共有リソース所有権

プロジェクト単位で一意な以下のリソースは **prod state でのみ create** する。dev / stg / qa state からは固定パスで参照するだけ。

| 共有リソース | 所有 state | 非 prod での扱い |
| --- | --- | --- |
| `google_artifact_registry_repository.server` (`server` Docker repo) | prod | image URI 文字列で参照のみ (`local.image`) |
| `google_project_iam_custom_role.fcm_sender` (`fcmSender` カスタムロール) | prod | `projects/<project>/roles/fcmSender` 固定パスで IAM binding |

この前提により、**非 prod env の apply 前に prod env の apply が完了している必要がある**。CI/CD では prod を先に通すパイプラインに固定すること。

## 初回ブートストラップ

Cloud Run Service / Job を立てる前に Artifact Registry と SA / Cloud SQL ユーザを先に作る必要がある (chicken-and-egg 回避)。`batch-jobs/infra` と同じ手順を踏む。

```bash
# 1. プロジェクト・GCS state バケット・Cloud SQL インスタンスは事前準備済み (本リポジトリの管理対象外)
# 2. Workload Identity Federation の設定 (way 配下のドキュメント参照)
# 3. envs/<env>.tfvars の REPLACE_ME を実 project_id / instance_connection_name に書き換える
#    (実値は社内 wiki / Secret Manager 参照。コミットしないこと)

# 4. AR / SA / IAM だけ先に apply (prod から)
terraform init -reconfigure -backend-config="prefix=server/prod"
terraform apply \
    -var-file=envs/prod.tfvars \
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

# 5. 1 度 Docker イメージを push (GitHub Actions deploy.yml で自動化される)

# 6. prod の残りを apply
terraform apply -var-file=envs/prod.tfvars -var=image_tag=<sha>

# 7. dev / stg / qa も同様に (共有リソースは create されず参照だけ)
terraform init -reconfigure -backend-config="prefix=server/dev"
terraform apply -var-file=envs/dev.tfvars -var=image_tag=<sha>
```

## 既存 academic-api Cloud Run Service の取り込み (cutover)

cutover 段階で `academic-api` / `academic-api-stg` / `academic-api-dev` Cloud Run Service を server state に取り込む。リソース名・リージョン・SA・環境変数キーを既存に合わせて差分ゼロにする。state は env ごとに分けてあるので、import も env ごとに実行する。

```bash
terraform init -reconfigure -backend-config="prefix=server/dev"
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
