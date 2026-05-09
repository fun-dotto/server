locals {
  # 環境ごとのリソース名サフィックス。way の set-env 規約と揃える:
  # - prod: <resource>
  # - dev / stg / qa: <resource>-<env>
  env_suffix = var.environment == "prod" ? "" : "-${var.environment}"

  # Artifact Registry repo は「server」固定 (環境跨ぎで 1 イメージを共有する)
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
  http_services = {
    "academic-api" = {
      command = ["/bin/academic-api"]
      cpu     = "1"
      memory  = "512Mi"
    }
  }

  # Cloud Run Job 定義。schedule = null の Job は Cloud Scheduler を作らない。
  cloud_run_jobs = {
    "build-class-change-notifications-job" = {
      command  = ["/bin/build-class-change-notifications-job"]
      schedule = var.build_class_change_notifications_schedule
      args     = []
      cpu      = "1"
      memory   = "512Mi"
      timeout  = "900s"
    }
    "dispatch-notifications-job" = {
      command  = ["/bin/dispatch-notifications-job"]
      schedule = var.dispatch_notifications_schedule
      args     = []
      cpu      = "1"
      memory   = "512Mi"
      timeout  = "900s"
    }
    "migrate-job" = {
      command  = ["/bin/migrate-job"]
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

  # SA は Service / Job ごとに 1 つずつ
  service_account_ids = concat(
    [for k, _ in local.http_services : k],
    [for k, _ in local.cloud_run_jobs : k],
  )
}
