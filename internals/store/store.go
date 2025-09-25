package store

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Metadata struct {
	Id           string
	ExportPath   string
	ExportMethod string
}

type Data interface {
	Metadata() Metadata
}

type DataStore[T Data] interface {
	/*Store a new chunk of data*/
	Store(data T)
	/*Get all data from store WITHOUT removing*/
	GetAll() []T
	/*Invalidate data from store, based on the provided metadata id*/
	Invalidate(id string)
}

type DummyData struct{}

func (d DummyData) Metadata() Metadata {
	return Metadata{
		Id:           uuid.New().String(),
		ExportPath:   "/health",
		ExportMethod: http.MethodGet,
	}
}

type DummyStore struct{}

func (*DummyStore) Store(data Data) {}
func (*DummyStore) GetAll() []Data {
	data := make([]Data, 0)
	data = append(data, DummyData{})

	return data
}
func (*DummyStore) Invalidate(id string) {
	fmt.Printf("Removing id %s from data store\n", id)
}
