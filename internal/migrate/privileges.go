// Package migrate に対してテーブル権限の再付与処理を提供する。
// shared-go の applyTablePrivileges を *sql.DB ベースに移植したもの。
package migrate

import (
	"context"
	"database/sql"
	"fmt"
)

// ApplyTablePrivileges は current_schema() 配下の全テーブルの所有者を
// dotto_admin に変更し、各ロールへ必要な権限を再付与する。
func ApplyTablePrivileges(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(
		ctx,
		"SELECT schemaname, tablename FROM pg_tables WHERE schemaname = current_schema()",
	)
	if err != nil {
		return fmt.Errorf("list tables: %w", err)
	}
	defer rows.Close()

	type table struct {
		schema string
		name   string
	}
	var tables []table
	for rows.Next() {
		var t table
		if err := rows.Scan(&t.schema, &t.name); err != nil {
			return fmt.Errorf("scan table row: %w", err)
		}
		tables = append(tables, t)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate table rows: %w", err)
	}

	for _, t := range tables {
		qualified := quoteIdent(t.schema) + "." + quoteIdent(t.name)
		stmts := []string{
			fmt.Sprintf("ALTER TABLE %s OWNER TO %s", qualified, quoteIdent("dotto_admin")),
			fmt.Sprintf("GRANT ALL PRIVILEGES ON TABLE %s TO %s", qualified, quoteIdent("dotto_admin")),
			fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE %s TO %s", qualified, quoteIdent("dotto_service")),
			fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE %s TO %s", qualified, quoteIdent("dotto_developer")),
		}
		for _, stmt := range stmts {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("exec %q: %w", stmt, err)
			}
		}
	}
	return nil
}
