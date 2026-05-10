data "external_schema" "gorm" {
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
  src = data.external_schema.gorm.url
  dev = "docker://postgres/18/dev"

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
