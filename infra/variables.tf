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

  validation {
    # google_sql_user は var.project_id 上で connection name の instance 部分のみを使うため、
    # connection name の project 部分が project_id と一致していないと
    # 「Cloud Run が接続する Cloud SQL」と「Terraform が IAM ユーザーを作る Cloud SQL」が
    # 別プロジェクトに分かれる事故を起こす。
    condition     = split(":", var.instance_connection_name)[0] == var.project_id
    error_message = "instance_connection_name の project 部分は var.project_id と一致させてください (別プロジェクトの Cloud SQL を参照する構成は未サポート)。"
  }
}

variable "db_name" {
  type        = string
  description = "接続先 PostgreSQL データベース名"
}

variable "image_tag" {
  type        = string
  description = "Cloud Run Service / Job が参照する Docker イメージタグ (commit SHA を想定)。prod は revision 追跡のため latest 禁止。"
  default     = "latest"

  validation {
    # 空文字や空白だけだと local.image が ".../server:" となり Cloud Run apply 時に
    # 不正な image URI で失敗する。
    condition     = trimspace(var.image_tag) != ""
    error_message = "image_tag に空文字や空白だけの値は指定できません。"
  }

  validation {
    # prod では追跡可能なタグ (commit SHA など) を必ず明示させる。
    # 他 env はローカル検証用に latest フォールバックを許容する。
    condition     = var.environment != "prod" || var.image_tag != "latest"
    error_message = "prod 環境では image_tag に \"latest\" を指定できません。commit SHA など追跡可能なタグを指定してください。"
  }
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
