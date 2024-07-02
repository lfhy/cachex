package cachex

import lru "github.com/hashicorp/golang-lru"

type LruCache[T any] struct {
	cache *lru.Cache
}

// Lru最近最少被使用的数据最先被淘汰
func NewMemLruCacheStroage(cacheSize ...int) *LruCache[any] {
	return NewMemLruCacheStroageWithType[any](cacheSize...)
}

// Lru最近最少被使用的数据最先被淘汰
func NewMemLruCacheStroageWithType[T any](cacheSize ...int) *LruCache[T] {
	cache := 1000000
	if len(cacheSize) > 0 && cacheSize[0] > 0 {
		cache = cacheSize[0]
	}
	lrucache, _ := lru.New(cache)
	return &LruCache[T]{
		cache: lrucache,
	}
}

// 获取缓存
func (c *LruCache[T]) Get(key string) (T, bool) {
	value, ok := c.cache.Get(key)
	if ok {
		return value.(T), true
	}
	var t T
	return t, false
}

// 设置缓存
func (c *LruCache[T]) Set(key string, value T) {
	c.cache.Add(key, value)
}

// 删除缓存
func (c *LruCache[T]) Delete(key string) {
	c.cache.Remove(key)
}

// 清空缓存
func (c *LruCache[T]) Free() {
	c.cache.Purge()
}
