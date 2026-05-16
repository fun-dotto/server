# 全バイナリ (academic-api / *-job / migrate-job) を 1 イメージに同梱するため
# Docker repo は server 1 つに集約する。環境ごとには分けない。
#
# プロジェクト単一の共有リソースなので、所有権は prod state に固定する。
# 他 env では既に作成済みのものを参照する形になるため、image URI は
# local.image (純粋な文字列) から組み立てて resource 参照には依存させない。
# 非 prod を apply する前に prod の AR repo apply が必須 (README 参照)。
resource "google_artifact_registry_repository" "server" {
  count = var.environment == "prod" ? 1 : 0

  location      = var.region
  repository_id = local.artifact_registry_repo
  format        = "DOCKER"
  description   = "Modular monolith server image (multi-binary)"

  depends_on = [google_project_service.required_apis]
}
