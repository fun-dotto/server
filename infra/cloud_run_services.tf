resource "google_cloud_run_v2_service" "http" {
  for_each = local.http_services

  name                = "${each.key}${local.env_suffix}"
  location            = var.region
  deletion_protection = var.environment == "prod"
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.workload[each.key].email
    # gen2 を明示。gen1 と比べてネットワーク / FS の挙動が安定し、
    # Direct VPC egress 等の将来要件にも gen2 が前提。
    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"

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
