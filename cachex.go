package cachex

type Cache[T any] struct {
	memcache    CacheStroage[T]
	remotecache []CacheStroage[T]
}

// 默认创any类型的存储
func NewCache(memcache CacheStroage[any], remotecache ...CacheStroage[any]) *Cache[any] {
	return &Cache[any]{
		memcache:    memcache,
		remotecache: remotecache,
	}
}

// 带有类型的存储
func NewCacheWithType[T any](memcache CacheStroage[T], remotecache ...CacheStroage[T]) *Cache[T] {
	return &Cache[T]{
		memcache:    memcache,
		remotecache: remotecache,
	}
}

// 获取缓存
func (c *Cache[T]) Get(key string) (T, bool) {
	value, ok := c.memcache.Get(key)
	if ok {
		return value, true
	}
	for _, remote := range c.remotecache {
		value, ok := remote.Get(key)
		if ok {
			c.memcache.Set(key, value)
			return value, true
		}
	}
	return value, false
}

// 设置缓存
func (c *Cache[T]) Set(key string, value T) {
	c.memcache.Set(key, value)
	for _, remote := range c.remotecache {
		remote.Set(key, value)
	}
}

// 删除缓存
func (c *Cache[T]) Delete(key string) {
	c.memcache.Delete(key)
	for _, remote := range c.remotecache {
		remote.Delete(key)
	}
}

// 获取失败时则设置缓存为指定值
func (c *Cache[T]) GetOrSet(key string, value T) (oldValue T, isLoad bool) {
	oldValue, ok := c.Get(key)
	if ok {
		return oldValue, true
	}
	c.Set(key, value)
	return value, false
}

// 获取失败时则通过指定函数设置为指定值
func (c *Cache[T]) GetOrSetFunc(key string, setFn func() (value T, canSet bool)) (oldValue T, isLoad bool) {
	oldValue, ok := c.Get(key)
	if ok {
		return oldValue, true
	}
	value, ok := setFn()
	if ok {
		c.Set(key, value)
		return value, false
	}
	return oldValue, false
}

// 获取缓存并设置为另一个值
func (c *Cache[T]) Swap(key string, value T) (oldValue T, isLoad bool) {
	oldValue, ok := c.Get(key)
	c.Set(key, value)
	return oldValue, ok
}

// 清空缓存
func (c *Cache[T]) Free() {
	c.memcache.Free()
	for _, remote := range c.remotecache {
		remote.Free()
	}
}
