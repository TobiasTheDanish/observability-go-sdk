package store

import "slices"

/*
A struct that implements the store.Store interface
with in-Memory data storage.
*/
type Memory struct {
	data []Data
}

func (m *Memory) Store(data Data) {
	if m.data == nil {
		m.data = make([]Data, 0)
	}

	m.data = append(m.data, data)
}

func (m *Memory) GetAll() []Data {
	if m.data == nil {
		m.data = make([]Data, 0)
	}

	return m.data
}

func (m *Memory) Invalidate(id string) {
	if m.data == nil {
		m.data = make([]Data, 0)
		return
	}

	idx := slices.IndexFunc(m.data, func(arg Data) bool {
		return arg.Metadata().Id == id
	})

	if idx == -1 {
		return
	}

	newData := make([]Data, 0)
	newData = append(newData, m.data[:idx]...)
	newData = append(newData, m.data[idx+1:]...)
	m.data = newData
}
