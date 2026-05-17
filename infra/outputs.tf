output "artifact_registry_repository" {
  value       = "${local.artifact_registry_host}/${var.project_id}/${local.artifact_registry_repo}"
  description = "Docker イメージの push 先リポジトリ"
}

output "image" {
  value       = local.image
  description = "Cloud Run Service / Job が共通で参照する Docker イメージ URI"
}

output "service_account_emails" {
  value       = { for k, sa in google_service_account.workload : k => sa.email }
  description = "Service / Job ごとの実行 SA email"
}

output "cloud_run_service_urls" {
  value       = { for k, s in google_cloud_run_v2_service.http : k => s.uri }
  description = "HTTP Cloud Run Service の URL マップ"
}

output "cloud_run_job_names" {
  value       = { for k, j in google_cloud_run_v2_job.jobs : k => j.name }
  description = "Cloud Run Job 名マップ (gcloud run jobs execute で利用)"
}

output "migrate_job_name" {
  value       = google_cloud_run_v2_job.jobs["migrate-job"].name
  description = "デプロイ後に実行する migrate Cloud Run Job 名"
}
