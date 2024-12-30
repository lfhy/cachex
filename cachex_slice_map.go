package cachex

import "sync"

type SliceMap[T any] struct {
	lock sync.RWMutex
	noCloseStroage
	data []*Data[T]
}

func (m *SliceMap[T]) Get(key string) (T, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, data := range m.data {
		if data.key == key {
			return data.value, true
		}
	}
	var t T
	return t, false
}

func (m *SliceMap[T]) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.delete(key)
}

// nolock
func (m *SliceMap[T]) delete(key string) {
	for index, data := range m.data {
		if data.key == key {
			m.data = append(m.data[:index], m.data[index+1:]...)
			return
		}
	}
}

// 获取第一个元素 并删除
func (m *SliceMap[T]) Pop() (T, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for index, data := range m.data {
		if index == 0 {
			m.data = append(m.data[:index], m.data[index+1:]...)
			return data.value, true
		}
	}
	var t T
	return t, false
}

func (m *SliceMap[T]) Set(key string, value T) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.delete(key)
	m.data = append(m.data, &Data[T]{
		key:   key,
		value: value,
	})
}

func (m *SliceMap[T]) Len() int {
	return len(m.data)
}

func (m *SliceMap[T]) Range(fn func(key string, value T) (Continue bool)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, data := range m.data {
		if !fn(data.key, data.value) {
			break
		}
	}
}

type Data[T any] struct {
	key   string
	value T
}
