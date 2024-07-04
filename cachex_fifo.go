package cachex

type FifoCache[T any] struct {
	cache *SliceMap[T]
	size  int
}

// Fifo最近最先加入缓存的数据最先被淘汰
func NewMemFifoCacheStroage(cacheSize ...int) *FifoCache[any] {
	return NewMemFifoCacheStroageWithType[any](cacheSize...)
}

// Fifo最近最先加入缓存的数据最先被淘汰
func NewMemFifoCacheStroageWithType[T any](cacheSize ...int) *FifoCache[T] {
	cache := 1000000
	if len(cacheSize) > 0 && cacheSize[0] > 0 {
		cache = cacheSize[0]
	}
	var data SliceMap[T]
	data.data = make([]*Data[T], 0)
	return &FifoCache[T]{
		size:  cache,
		cache: &data,
	}
}

// 获取缓存
func (c *FifoCache[T]) Get(key string) (T, bool) {
	value, ok := c.cache.Get(key)
	if ok {
		return value, true
	}
	var t T
	return t, false
}

// 设置缓存
func (c *FifoCache[T]) Set(key string, value T) {
	if c.cache.Len() >= c.size {
		c.cache.Pop()
	}
	c.cache.Set(key, value)
}

// 删除缓存
func (c *FifoCache[T]) Delete(key string) {
	c.cache.Delete(key)
}

// 清空缓存
func (c *FifoCache[T]) Free() {
	cache := SliceMap[T]{
		data: make([]*Data[T], 0),
	}
	c.cache = &cache
}
