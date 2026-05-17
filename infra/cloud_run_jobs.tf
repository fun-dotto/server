resource "google_cloud_run_v2_job" "jobs" {
  for_each = local.cloud_run_jobs

  name                = "${each.key}${local.env_suffix}"
  location            = var.region
  deletion_protection = var.environment == "prod"

  template {
    template {
      service_account = google_service_account.workload[each.key].email
      timeout         = each.value.timeout
      max_retries     = each.value.max_retries
      # gen2 を明示 (cloud_run_services.tf と同じ理由)。
      execution_environment = "EXECUTION_ENVIRONMENT_GEN2"

      containers {
        image   = local.image
        command = each.value.command
        args    = each.value.args

        dynamic "env" {
          for_each = local.common_db_env
          content {
            name  = env.key
            value = env.value
          }
        }

        env {
          name  = "DB_IAM_USER"
          value = trimsuffix(google_service_account.workload[each.key].email, ".gserviceaccount.com")
        }

        # Firebase Admin SDK が project ID を ADC から推測するが、
        # Cloud Run Jobs では明示しておくと安全 (dispatch-notifications-job 用)。
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.project_id
        }

        resources {
          limits = {
            cpu    = each.value.cpu
            memory = each.value.memory
          }
        }
      }
    }
  }

  depends_on = [
    google_project_service.required_apis,
    google_artifact_registry_repository.server,
    google_project_iam_member.workload_sql_client,
    google_project_iam_member.workload_sql_instance_user,
    google_sql_user.workload,
  ]
}
