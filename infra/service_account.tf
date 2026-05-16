# Service / Job ごとに専用 SA を 1 つずつ作成 (最小権限の原則を Phase 1 から徹底)。
# - map key = バイナリ名 (Service / Job 側から google_service_account.workload[<binary>] で参照)
# - account_id = 短縮した sa_id + env_suffix (30 字制限内)
# - display_name はコンソール識別性のためフルのバイナリ名を採用
resource "google_service_account" "workload" {
  for_each = local.service_account_sa_ids

  account_id   = "${each.value}${local.env_suffix}"
  display_name = "${each.key}${local.env_suffix}"

  depends_on = [google_project_service.required_apis]
}

# 各 SA に Cloud SQL クライアント権限と IAM ユーザ権限を付与
resource "google_project_iam_member" "workload_sql_client" {
  for_each = google_service_account.workload

  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${each.value.email}"
}

resource "google_project_iam_member" "workload_sql_instance_user" {
  for_each = google_service_account.workload

  project = var.project_id
  role    = "roles/cloudsql.instanceUser"
  member  = "serviceAccount:${each.value.email}"
}

# 各 SA を Cloud SQL の IAM ユーザーとして登録
resource "google_sql_user" "workload" {
  for_each = google_service_account.workload

  name     = trimsuffix(each.value.email, ".gserviceaccount.com")
  instance = local.cloud_sql_instance_name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"

  depends_on = [google_project_service.required_apis]
}

# dispatch-notifications-job が FCM HTTP v1 API でプッシュ通知を送るための最小権限。
# server 側を「正」として fcm_sender カスタムロールを定義し (Notion 計画 §4 / §9)、
# batch-jobs 側は cutover 時に data リソースで参照のみへ切り替える。
#
# プロジェクト単一なので、本リポジトリ内では prod state のみがこのカスタムロールを
# 所有する (count で制御)。非 prod env を apply する前に prod を先行 apply する運用。
resource "google_project_iam_custom_role" "fcm_sender" {
  count = var.environment == "prod" ? 1 : 0

  role_id     = "fcmSender"
  title       = "FCM Sender"
  description = "Send FCM messages via the HTTP v1 API"
  permissions = [
    "cloudmessaging.messages.create",
  ]
}

# fcm_sender ロールへの付与は env ごとに行う (dispatch-notifications-job の SA が env 別なので)。
# 非 prod では prod state が作成済みのロールを固定パスで参照する。
resource "google_project_iam_member" "dispatch_notifications_fcm_sender" {
  project = var.project_id
  role    = local.fcm_sender_role
  member  = "serviceAccount:${google_service_account.workload["dispatch-notifications-job"].email}"

  depends_on = [google_project_iam_custom_role.fcm_sender]
}

# Cloud Scheduler 用 SA。Cloud Run Job (googleapis.com) 起動時の OAuth access token 発行に使う。
resource "google_service_account" "scheduler" {
  account_id   = "scheduler${local.env_suffix}"
  display_name = "Cloud Scheduler invoker${local.env_suffix}"

  depends_on = [google_project_service.required_apis]
}

resource "google_service_account_iam_member" "scheduler_token_creator" {
  service_account_id = google_service_account.scheduler.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:service-${data.google_project.current.number}@gcp-sa-cloudscheduler.iam.gserviceaccount.com"

  depends_on = [google_project_service.required_apis]
}
