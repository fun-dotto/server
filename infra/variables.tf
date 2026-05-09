variable "project_id" {
  type        = string
  description = "Google Cloud プロジェクト ID (全環境を 1 プロジェクトに同居させる)"
}

variable "region" {
  type        = string
  description = "Cloud Run Service / Job / Artifact Registry / Scheduler のリージョン"
  default     = "asia-northeast1"
}

variable "environment" {
  type        = string
  description = "デプロイ環境 (dev / stg / qa / prod)。命名規約: prod は <resource>、それ以外は <resource>-<env>"

  validation {
    condition     = contains(["dev", "stg", "qa", "prod"], var.environment)
    error_message = "environment は dev / stg / qa / prod のいずれかを指定してください。"
  }
}

variable "instance_connection_name" {
  type        = string
  description = "Cloud SQL のインスタンス接続名 (project:region:instance)"

  validation {
    condition = length(split(":", var.instance_connection_name)) == 3 && alltrue([
      for part in split(":", var.instance_connection_name) : trimspace(part) != ""
    ])
    error_message = "instance_connection_name は \"project:region:instance\" 形式で指定してください。"
  }
}

variable "db_name" {
  type        = string
  description = "接続先 PostgreSQL データベース名"
}

variable "image_tag" {
  type        = string
  description = "Cloud Run Service / Job が参照する Docker イメージタグ (commit SHA を想定)"
  default     = "latest"
}

variable "secret_project_id" {
  type        = string
  description = "Secret Manager を保持するプロジェクト ID。Phase 1 では未使用だが命名と紐付けのため宣言だけ残す"
  default     = ""
}

variable "build_class_change_notifications_schedule" {
  type        = string
  description = "build-class-change-notifications-job の cron 式 (Asia/Tokyo)"
  default     = "30 17 * * *"
}

variable "dispatch_notifications_schedule" {
  type        = string
  description = "dispatch-notifications-job の cron 式 (Asia/Tokyo)"
  default     = "0 18 * * *"
}
