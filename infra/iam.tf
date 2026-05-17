# HTTP サービス (academic-api ほか) は Phase 1 では allUsers での公開を許可する。
# Phase 4 で BFF を再評価する際に内部公開などへ絞り直す。
resource "google_cloud_run_v2_service_iam_member" "http_public_invoker" {
  for_each = google_cloud_run_v2_service.http

  project  = each.value.project
  location = each.value.location
  name     = each.value.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
