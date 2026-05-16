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

  # cleanup_policies を実行する (true にすると評価のみで削除されない)。
  cleanup_policy_dry_run = false

  # commit SHA 単位で日々イメージが増える前提のため、最初から retention を入れる。
  # - 直近 50 タグはロールバック余地として残す
  # - タグ未付与の untagged manifest は 7 日で削除 (CI 中間レイヤなど)
  cleanup_policies {
    id     = "keep-recent-tagged"
    action = "KEEP"
    most_recent_versions {
      keep_count = 50
    }
  }

  cleanup_policies {
    id     = "delete-untagged-after-7d"
    action = "DELETE"
    condition {
      tag_state  = "UNTAGGED"
      older_than = "604800s"
    }
  }

  # 全 env が文字列パスで暗黙参照しているため、prod state で誤って destroy
  # されると非 prod の Cloud Run image pull が即死する。
  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.required_apis]
}
