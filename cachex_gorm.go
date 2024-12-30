package cachex

import "gorm.io/gorm"

type CachexGormModel struct {
	tableName string
	ID        int64  `gorm:"id"`
	Key       string `gorm:"key"`
	Value     string `gorm:"value"`
}

func (c CachexGormModel) TableName() string {
	if c.tableName != "" {
		return c.tableName
	}
	return "cachex"
}

type CachexGorm[T any] struct {
	tableName string
	db        *gorm.DB
}

func NewGormCachex[T any](db *gorm.DB, tableName ...string) *CachexGorm[T] {
	table := "cachex"
	if len(tableName) > 0 {
		table = tableName[0]
	}
	db.AutoMigrate(&CachexGormModel{tableName: table})
	return &CachexGorm[T]{db: db, tableName: table}
}

func (c *CachexGorm[T]) NewModel() CachexGormModel {
	return CachexGormModel{tableName: c.tableName}
}

func (c *CachexGorm[T]) Get(key string) (T, bool) {
	var res T
	value := c.NewModel()
	err := c.db.Where("key = ?", key).First(&value).Error
	if err != nil {
		return res, false
	}
	return res, unmarshalString(value.Value, &res)
}

func (c *CachexGorm[T]) Set(key string, value T) {
	if key == "" {
		return
	}
	set := c.NewModel()
	set.Key = key
	set.Value = marshalString(value)
	before := c.NewModel()
	err := c.db.Model(&before).Where("key = ?", key).First(&before).Error
	if err != nil {
		c.db.Create(&set)
	} else {
		c.db.Model(&before).Updates(set)
	}
}

func (c *CachexGorm[T]) Delete(key string) {
	if key == "" {
		return
	}
	value := c.NewModel()
	value.Key = key
	c.db.Model(value).Where("key = ?", key).Delete(value)
}

func (c *CachexGorm[T]) Free() {
	c.db.Delete(c.NewModel())
}

func (c *CachexGorm[T]) GetDBInterface() *gorm.DB {
	return c.db
}

func (c *CachexGorm[T]) Close() {
}

func (c *CachexGorm[T]) Range(f func(key string, value T) bool) {
	res, err := c.db.Model(&CachexGormModel{}).Find(&[]CachexGormModel{}).Rows()
	if err != nil {
		return
	}
	for res.Next() {
		value := c.NewModel()
		res.Scan(&value)
		var t T
		if unmarshalString(value.Value, &t) {
			if !f(value.Key, t) {
				break
			}
		}
	}
}
