package hashmap

import "encoding/json"

type Map[K comparable, V any] struct {
	keys   []K
	values []V
	data   map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		keys:   []K{},
		values: []V{},
		data:   make(map[K]V),
	}
}

func (m *Map[K, V]) Put(key K, value V) {
	if _, exists := m.data[key]; !exists {
		m.keys = append(m.keys, key)
		m.values = append(m.values, value)
	}
	m.data[key] = value
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	value, exists := m.data[key]
	return value, exists
}

func (m *Map[K, V]) Delete(key K) {
	if _, exists := m.data[key]; exists {
		delete(m.data, key)
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				m.values = append(m.values[:i], m.values[i+1:]...)
				break
			}
		}
	}
}

func (m *Map[K, V]) Contains(key K) bool {
	_, exists := m.data[key]
	return exists
}

func (m *Map[K, V]) Keys() []K {
	return m.keys
}

func (m *Map[K, V]) Values() []V {
	return m.values
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	tempMap := make(map[K]V)

	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}
	m.data = make(map[K]V)
	m.keys = make([]K, 0, len(tempMap))
	m.values = make([]V, 0, len(tempMap))
	for key, value := range tempMap {
		m.keys = append(m.keys, key)
		m.values = append(m.values, value)
		m.data[key] = value
	}
	return nil
}

func (m *Map[K, V]) Iter() <-chan struct {
	Key   K
	Value V
} {
	ch := make(chan struct {
		Key   K
		Value V
	})
	go func() {
		defer close(ch)
		for key, value := range m.data {
			ch <- struct {
				Key   K
				Value V
			}{Key: key, Value: value}
		}
	}()
	return ch
}
