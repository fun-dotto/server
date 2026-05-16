locals {
  # 環境ごとのリソース名サフィックス。way の set-env 規約と揃える:
  # - prod: <resource>
  # - dev / stg / qa: <resource>-<env>
  env_suffix = var.environment == "prod" ? "" : "-${var.environment}"

  # Artifact Registry repo は「server」固定 (環境跨ぎで 1 イメージを共有する)。
  # 所有権は prod state に固定する (artifact_registry.tf 参照)。
  artifact_registry_repo = "server"
  artifact_registry_host = "${var.region}-docker.pkg.dev"

  # 全バイナリを 1 イメージに同梱するため image は 1 本のみ。
  # 各 Cloud Run Service / Job は command で /bin/<name> を切り替える。
  image = "${local.artifact_registry_host}/${var.project_id}/${local.artifact_registry_repo}/${local.artifact_registry_repo}:${var.image_tag}"

  cloud_sql_instance_name = split(":", var.instance_connection_name)[2]

  # 共通 DB 環境変数。各 Service / Job の env として展開する。
  common_db_env = {
    INSTANCE_CONNECTION_NAME = var.instance_connection_name
    DB_NAME                  = var.db_name
  }

  # HTTP サービス定義。Phase 1 では academic-api のみ。Phase 2 で
  # announcement-api / user-api 等を追加する想定。
  #
  # sa_id: google_service_account.account_id は 30 文字までという GCP 制約があるため、
  #        SA 用に短縮した識別子を別途持たせる (env_suffix を足しても 30 字以内に収まること)。
  http_services = {
    "academic-api" = {
      sa_id   = "academic-api"
      command = ["/bin/academic-api"]
      cpu     = "1"
      memory  = "512Mi"
    }
  }

  # Cloud Run Job 定義。schedule = null の Job は Cloud Scheduler を作らない。
  cloud_run_jobs = {
    "build-class-change-notifications-job" = {
      sa_id    = "class-change-notif-job"
      command  = ["/bin/build-class-change-notifications-job"]
      schedule = var.build_class_change_notifications_schedule
      args     = []
      cpu      = "1"
      memory   = "512Mi"
      timeout  = "900s"
    }
    "dispatch-notifications-job" = {
      sa_id    = "dispatch-notif-job"
      command  = ["/bin/dispatch-notifications-job"]
      schedule = var.dispatch_notifications_schedule
      args     = []
      cpu      = "1"
      memory   = "512Mi"
      timeout  = "900s"
    }
    # Cloud Run Job 名は "migrate-job" (way 規約上 Job には -job 接尾辞) だが、
    # 同梱バイナリは cmd/migrate に対応する /bin/migrate を起動する。
    "migrate-job" = {
      sa_id    = "migrate-job"
      command  = ["/bin/migrate"]
      schedule = null
      args     = []
      cpu      = "1"
      memory   = "512Mi"
      timeout  = "900s"
    }
  }

  # Cloud Scheduler を作る Job だけ抽出
  scheduled_jobs = {
    for k, v in local.cloud_run_jobs : k => v if v.schedule != null
  }

  # SA は Service / Job ごとに 1 つずつ。map の key は他リソースから参照しやすい
  # バイナリ名 (cloud_run_jobs / http_services の key と一致)、value は SA の account_id。
  service_account_sa_ids = merge(
    { for k, v in local.http_services : k => v.sa_id },
    { for k, v in local.cloud_run_jobs : k => v.sa_id },
  )

  # 共有プロジェクト IAM カスタムロール (fcm_sender) は prod state でのみ create し、
  # 他 env からは projects/<project>/roles/<role_id> の固定パスで参照する。
  fcm_sender_role = "projects/${var.project_id}/roles/fcmSender"
}
