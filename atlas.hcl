data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas-loader",
  ]
}

# マイグレーション生成用（接続先は持たない）
env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/18/dev"
  migration {
    dir = "file://migrations"
  }
  lint {
    # 直近 1 ファイル（今回追加したマイグレーション）のみを lint 対象にする
    latest = 1
  }
}

# 適用先の環境。接続先は環境変数 DATABASE_URL で注入する。
# 例: postgres://<iam-user>@127.0.0.1:5432/<dotto | dotto_stg | dotto_dev>?sslmode=disable
# （Cloud SQL Auth Proxy --auto-iam-authn 経由。GitHub Actions では .github/actions/atlas-migrate が設定する）
env "dev" {
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}

env "stg" {
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}

env "prd" {
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}
