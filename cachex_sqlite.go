package cachex

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSqliteCache[T any](dbPath string, tableName ...string) (*CachexGorm[T], error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	})
	if err != nil {
		return nil, err
	}
	return NewGormCachex[T](db, tableName...), nil
}
