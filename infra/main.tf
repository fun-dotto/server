terraform {
  # variables.tf の image_tag validation で var.environment を参照するため、
  # cross-variable validation が GA になった Terraform 1.9 系以上を必須にする。
  required_version = ">= 1.9.0"

  # prefix は env ごとに分離するため partial configuration にしている。
  # `terraform init -backend-config="prefix=server/<env>"` が必須。
  # ここで prefix を固定しておくと 4 env で同一 state を共有する事故が起きるため、
  # 意図的に未指定 (-backend-config 忘れ時は init がエラーになる)。
  backend "gcs" {
    bucket = "swift2023groupc-tfstate"
  }

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
