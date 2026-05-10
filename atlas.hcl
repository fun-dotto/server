// Atlas プロジェクト設定。
// - 望ましいスキーマ (desired state) は internal/shared/model/ の GORM モデルを正とする。
// - migrations/ 配下に versioned SQL を保存し、サム (atlas.sum) で改ざん検知。
// - dev DB は使い捨ての Docker コンテナ (Postgres 18) を利用する。
// - cmd/migrate-job が ariga.io/atlas-go-sdk (atlasexec) 経由で migrations/ を適用する。

variable "url" {
  type    = string
  default = ""
}

data "external_schema" "app" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/shared/model",
    "--dialect", "postgres",
  ]
}

env "local" {
  src = data.external_schema.app.url
  dev = "docker://postgres/18/dev?search_path=public"

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  src = data.external_schema.app.url
  url = var.url
  dev = "docker://postgres/18/dev?search_path=public"

  migration {
    dir = "file://migrations"
  }
}
