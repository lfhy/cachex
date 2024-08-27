package cachex

type SliceMap[T any] struct {
	noCloseStroage
	data []*Data[T]
}

func (m *SliceMap[T]) Get(key string) (T, bool) {
	for _, data := range m.data {
		if data.key == key {
			return data.value, true
		}
	}
	var t T
	return t, false
}

func (m *SliceMap[T]) Delete(key string) {
	for index, data := range m.data {
		if data.key == key {
			m.data = append(m.data[:index], m.data[index+1:]...)
			return
		}
	}
}

func (m *SliceMap[T]) Pop() (T, bool) {
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
	m.Delete(key)
	m.data = append(m.data, &Data[T]{
		key:   key,
		value: value,
	})
}

func (m *SliceMap[T]) Len() int {
	return len(m.data)
}

type Data[T any] struct {
	key   string
	value T
}
