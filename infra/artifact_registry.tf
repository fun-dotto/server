# 全バイナリ (academic-api / *-job / migrate-job) を 1 イメージに同梱するため
# Docker repo は server 1 つに集約する。環境ごとには分けない。
resource "google_artifact_registry_repository" "server" {
  location      = var.region
  repository_id = local.artifact_registry_repo
  format        = "DOCKER"
  description   = "Modular monolith server image (multi-binary)"

  depends_on = [google_project_service.required_apis]
}
