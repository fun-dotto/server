package a

import (
	"gorm.io/gorm"
	"uuid"
)

func f(db *gorm.DB, id uuid.UUID, ids []uuid.UUID, args []any) {
	db.Where("id = ?", id)                  // want `GORM に uuid.UUID を直接渡さず`
	db.Where("id IN ?", ids)                // want `GORM に uuid.UUID を直接渡さず`
	db.First(&struct{}{}, "id = ?", &id)    // want `GORM に uuid.UUID を直接渡さず`
	db.Exec("INSERT INTO t VALUES (?)", id) // want `GORM に uuid.UUID を直接渡さず`
	db.Update("id", id)                     // want `GORM に uuid.UUID を直接渡さず`

	db.Where("id = ?", id.String())
	db.Where("id = ?", args...)
	db.First(&struct{}{}, "id = ?", "x")
}
