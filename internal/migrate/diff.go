package migrate

import (
	"context"
	"database/sql"

	"ariga.io/atlas-provider-gorm/gormschema"
	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
	"github.com/fun-dotto/server/internal/shared/model"
)

// GenerateDiff はGORMモデルとDBの差分を migration ファイルとして書き出します
func GenerateDiff(ctx context.Context, db *sql.DB, dirPath string, name string) error {
	// 1. Atlasドライバの準備
	drv, err := postgres.Open(db)
	if err != nil {
		return err
	}

	// 2. GORMモデルから理想のスキーマ定義をロード
	// 第2引数以降にすべてのGORMモデル構造体を渡します
	loader := gormschema.New("postgres")
	desiredSql, err := loader.Load(&model.Announcement{}, &model.User{})
	if err != nil {
		return err
	}

	// 3. 出力先ディレクトリの準備
	dir, err := migrate.NewLocalDir(dirPath)
	if err != nil {
		return err
	}

	// 4. 差分計算とファイル出力
	if err := drv.Diff(ctx, dir, desiredSql); err != nil {
		return err
	}

	// 5. ハッシュファイルの更新 (整合性チェック用)
	return dir.WriteHash()
}
