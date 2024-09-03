package cachex

import (
	"strings"

	leveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LeveldbCache[T any] struct {
	cache *leveldb.DB
}

// Leveldb最近最少被使用的数据最先被淘汰
func NewLevelDBCacheStroage(dbpath string) (*LeveldbCache[any], error) {
	return NewLevelDBCacheStroageWithType[any](dbpath)
}

// Leveldb最近最少被使用的数据最先被淘汰
func NewLevelDBCacheStroageWithType[T any](dbpath string) (*LeveldbCache[T], error) {
	// 打开或创建一个LevelDB数据库
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		if strings.Contains(err.Error(), "temporarily") {
			db, err = leveldb.RecoverFile(dbpath, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return &LeveldbCache[T]{
		cache: db,
	}, nil
}

// 获取缓存
func (c *LeveldbCache[T]) Get(key string) (T, bool) {
	var getData JsonByteData[T]
	data, err := c.cache.Get([]byte(key), &opt.ReadOptions{})
	if err != nil {
		return getData.Data, false
	}
	return getData.Data, unmarshal(data, &getData)
}

// 设置缓存
func (c *LeveldbCache[T]) Set(key string, value T) {
	var setData JsonByteData[T]
	setData.Data = value
	data := marshal(setData)
	c.cache.Put([]byte(key), data, &opt.WriteOptions{})
}

// 删除缓存
func (c *LeveldbCache[T]) Delete(key string) {
	c.cache.Delete([]byte(key), &opt.WriteOptions{})
}

// 清空缓存
func (c *LeveldbCache[T]) Free() {
	iter := c.cache.NewIterator(&util.Range{}, &opt.ReadOptions{})
	for iter.Next() {
		iter.Release()
	}
}

// 获取DB实例对象
func (c *LeveldbCache[T]) GetDBInterface() *leveldb.DB {
	return c.cache
}

// 获取DB实例对象
func (c *LeveldbCache[T]) Close() {
	c.cache.Close()
}
