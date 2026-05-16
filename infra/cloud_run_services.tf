resource "google_cloud_run_v2_service" "http" {
  for_each = local.http_services

  name                = "${each.key}${local.env_suffix}"
  location            = var.region
  deletion_protection = var.environment == "prod"
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.workload[each.key].email

    scaling {
      max_instance_count = each.value.max_instance_count
    }

    containers {
      image   = local.image
      command = each.value.command

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

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = each.value.cpu
          memory = each.value.memory
        }
      }
    }
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }

  depends_on = [
    google_project_service.required_apis,
    google_artifact_registry_repository.server,
    google_project_iam_member.workload_sql_client,
    google_project_iam_member.workload_sql_instance_user,
    google_sql_user.workload,
  ]
}
