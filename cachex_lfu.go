package cachex

import (
	"container/list"
	"sync"
)

type LfuCache[T any] struct {
	noCloseStroage
	UpperBound      int
	LowerBound      int
	values          map[string]*cacheEntry[T]
	freqs           *list.List
	len             int
	lock            sync.RWMutex
	EvictionChannel chan<- Eviction[T]
}

type cacheEntry[T any] struct {
	key      string
	value    T
	freqNode *list.Element
}

type listEntry[T any] struct {
	entries map[*cacheEntry[T]]byte
	freq    int
}

type Eviction[T any] struct {
	Key   string
	Value T
}

// Lfu使用频率最低的数据最先被淘汰
func NewMemLfuCacheStroage() *LfuCache[any] {
	return NewMemLfuCacheStroageWithType[any]()
}

// Lfu使用频率最低的数据最先被淘汰
func NewMemLfuCacheStroageWithType[T any]() *LfuCache[T] {
	return &LfuCache[T]{
		values: make(map[string]*cacheEntry[T]),
		freqs:  list.New(),
	}
}

// 获取缓存
func (c *LfuCache[T]) Get(key string) (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if e, ok := c.values[key]; ok {
		c.increment(e)
		return e.value, true
	}
	var t T
	return t, false
}

// 设置缓存
func (c *LfuCache[T]) Set(key string, value T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		// value already exists for key.  overwrite
		e.value = value
		c.increment(e)
	} else {
		// value doesn't exist.  insert
		var e cacheEntry[T]
		e.key = key
		e.value = value
		c.values[key] = &e
		c.increment(&e)
		c.len++
		// bounds mgmt
		if c.UpperBound > 0 && c.LowerBound > 0 {
			if c.len > c.UpperBound {
				c.evict(c.len - c.LowerBound)
			}
		}
	}
}

// 删除缓存
func (c *LfuCache[T]) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		c.remEntry(e.freqNode, e)
		delete(c.values, key)
	}
}

// 清空缓存
func (c *LfuCache[T]) Free() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.evict(c.len)
}

func (c *LfuCache[T]) evict(count int) int {
	// No lock here so it can be called
	// from within the lock (during Set)
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry := range place.Value.(*listEntry[T]).entries {
				if i < count {
					if c.EvictionChannel != nil {
						c.EvictionChannel <- Eviction[T]{
							Key:   entry.key,
							Value: entry.value,
						}
					}
					delete(c.values, entry.key)
					c.remEntry(place, entry)
					evicted++
					c.len--
					i++
				}
			}
		}
	}
	return evicted
}

func (c *LfuCache[T]) increment(e *cacheEntry[T]) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element
	if currentPlace == nil {
		// new entry
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		// move up
		nextFreq = currentPlace.Value.(*listEntry[T]).freq + 1
		nextPlace = currentPlace.Next()
	}

	if nextPlace == nil || nextPlace.Value.(*listEntry[T]).freq != nextFreq {
		// create a new list entry
		li := new(listEntry[T])
		li.freq = nextFreq
		li.entries = make(map[*cacheEntry[T]]byte)
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	nextPlace.Value.(*listEntry[T]).entries[e] = 1
	if currentPlace != nil {
		// remove from current position
		c.remEntry(currentPlace, e)
	}
}

func (c *LfuCache[T]) remEntry(place *list.Element, entry *cacheEntry[T]) {
	entries := place.Value.(*listEntry[T]).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}
