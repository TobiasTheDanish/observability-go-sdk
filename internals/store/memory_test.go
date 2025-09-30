package store

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type mockData struct {
	id string
}

func (d *mockData) Metadata() Metadata {
	if d.id == "" {
		d.id = uuid.NewString()
	}

	return Metadata{
		Id:           d.id,
		ExportPath:   "/",
		ExportMethod: "POST",
	}
}

func TestMemoryStore(t *testing.T) {
	tests := []struct {
		name     string
		op       func(store *Memory) *Memory
		expected []Data
	}{
		{
			name: "Single store operation",
			op: func(store *Memory) *Memory {
				data := &mockData{
					id: "mock-id",
				}

				store.Store(data)

				return store
			},
			expected: []Data{
				&mockData{
					id: "mock-id",
				},
			},
		},
		{
			name: "Single store operation with invalidation",
			op: func(store *Memory) *Memory {
				data := &mockData{
					id: "mock-id",
				}

				store.Store(data)
				store.Invalidate("mock-id")

				return store
			},
			expected: []Data{},
		},
	}

	for _, tt := range tests {
		store := &Memory{}
		t.Run(tt.name, func(t *testing.T) {
			store = tt.op(store)

			if !reflect.DeepEqual(store.data, tt.expected) {
				t.Fatalf("Memory data slice did not match expected. Got: %v, expected: %v", store.data, tt.expected)
			}
		})
	}
}
