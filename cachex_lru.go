package cachex

import lru "github.com/karlseguin/ccache/v3"

type LruCache[T any] struct {
	noCloseStroage
	cache *lru.Cache[T]
}

// Lru最近最少被使用的数据最先被淘汰
func NewMemLruCacheStroage(cacheSize ...int) *LruCache[any] {
	return NewMemLruCacheStroageWithType[any](cacheSize...)
}

// Lru最近最少被使用的数据最先被淘汰
func NewMemLruCacheStroageWithType[T any](cacheSize ...int) *LruCache[T] {
	cache := 50000
	if len(cacheSize) > 0 && cacheSize[0] > 0 {
		cache = cacheSize[0]
	}
	lrucache := lru.New(lru.Configure[T]().MaxSize(int64(cache)))
	return &LruCache[T]{
		cache: lrucache,
	}
}

// 获取缓存
func (c *LruCache[T]) Get(key string) (T, bool) {
	value := c.cache.Get(key)
	if value != nil {
		return value.Value(), true
	}
	var t T
	return t, false
}

// 设置缓存
func (c *LruCache[T]) Set(key string, value T) {
	c.cache.Set(key, value, 0)
}

// 删除缓存
func (c *LruCache[T]) Delete(key string) {
	c.cache.Delete(key)
}

// 清空缓存
func (c *LruCache[T]) Free() {
	c.cache.Clear()
}

// 遍历数据
func (c *LruCache[T]) Range(f func(key string, value T) bool) {
	c.cache.ForEachFunc(func(key string, item *lru.Item[T]) bool {
		return f(key, item.Value())
	})
}

// 获取DB实例对象
func (c *LruCache[T]) GetDBInterface() *lru.Cache[T] {
	return c.cache
}

// 获取全部缓存Key
func (c *LruCache[T]) Keys() []string {
	var keys []string
	c.cache.ForEachFunc(func(key string, item *lru.Item[T]) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}
