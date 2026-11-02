package logging

import (
	"log/slog"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

// GormLogger は SQL エラーとスロークエリを slog に出力する GORM ロガー。
// record not found は呼び出し側で 404 として扱うため出力しない。
func GormLogger() gormlogger.Interface {
	return gormlogger.NewSlogLogger(slog.Default(), gormlogger.Config{
		SlowThreshold:             500 * time.Millisecond,
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: true,
	})
}
