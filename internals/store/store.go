package store

type DataStore[T any] interface {
	Store(data T)
	GetAll() []T
}

type DummyStore struct{}

func (*DummyStore) Store(data any) {}
func (*DummyStore) GetAll() []any {
	return make([]any, 1, 1)
}
