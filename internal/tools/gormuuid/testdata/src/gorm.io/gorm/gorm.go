package gorm

type DB struct{}

func (db *DB) Where(query any, args ...any) *DB    { return db }
func (db *DB) First(dest any, conds ...any) *DB    { return db }
func (db *DB) Exec(sql string, values ...any) *DB  { return db }
func (db *DB) Update(column string, value any) *DB { return db }
