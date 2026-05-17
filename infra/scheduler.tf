resource "google_cloud_run_v2_job_iam_member" "scheduler_invoker" {
  for_each = local.scheduled_jobs

  project  = google_cloud_run_v2_job.jobs[each.key].project
  location = google_cloud_run_v2_job.jobs[each.key].location
  name     = google_cloud_run_v2_job.jobs[each.key].name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.scheduler.email}"
}

resource "google_cloud_scheduler_job" "triggers" {
  for_each = local.scheduled_jobs

  name        = "${each.key}-trigger${local.env_suffix}"
  description = "Trigger ${each.key}${local.env_suffix} Cloud Run Job"
  schedule    = each.value.schedule
  time_zone   = "Asia/Tokyo"
  region      = var.region

  http_target {
    http_method = "POST"
    # Job リソース属性から URI を組み立てて、google_cloud_run_v2_job.jobs 側で
    # name / location / project の派生方法が変わってもトリガー先がずれないようにする。
    uri = "https://run.googleapis.com/v2/projects/${google_cloud_run_v2_job.jobs[each.key].project}/locations/${google_cloud_run_v2_job.jobs[each.key].location}/jobs/${google_cloud_run_v2_job.jobs[each.key].name}:run"
    headers = {
      "Content-Type" = "application/json"
    }
    body = base64encode("{}")

    oauth_token {
      service_account_email = google_service_account.scheduler.email
    }
  }

  depends_on = [
    google_project_service.required_apis,
    google_cloud_run_v2_job_iam_member.scheduler_invoker,
  ]
}
