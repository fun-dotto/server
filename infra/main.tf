terraform {
  # variables.tf の image_tag validation で var.environment を参照するため、
  # cross-variable validation が GA になった Terraform 1.9 系以上を必須にする。
  required_version = ">= 1.9.0"

  # 4 env で同一 state を共有する事故を防ぐため、bucket / prefix とも
  # ここでは固定せず envs/<env>.backend.hcl 経由で渡す partial configuration にする。
  # bucket を未指定にすることで素の `terraform init` は init エラーで弾かれ、
  # envs/<env>.backend.hcl を渡さない限り bucket 直下のデフォルト state を
  # 誤って使ってしまう経路が無くなる。
  # 例: `terraform init -reconfigure -backend-config=envs/dev.backend.hcl`
  backend "gcs" {}

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "required_apis" {
  for_each = toset([
    "run.googleapis.com",
    "cloudscheduler.googleapis.com",
    "secretmanager.googleapis.com",
    "artifactregistry.googleapis.com",
    "sqladmin.googleapis.com",
    "iam.googleapis.com",
    "fcm.googleapis.com",
    "firebase.googleapis.com",
  ])

  service            = each.key
  disable_on_destroy = false
}

data "google_project" "current" {
  project_id = var.project_id
}
