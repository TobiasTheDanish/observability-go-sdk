package export

import (
	"fmt"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/tobiasthedanish/observability-go-sdk/internals/store"
)

var (
	metaUuid string = "f68ff9fb-9082-4e8c-854a-6251e93194be"
)

type MockData struct{}

func (d MockData) Metadata() store.Metadata {
	return store.Metadata{
		Id:           metaUuid,
		ExportPath:   "/health",
		ExportMethod: http.MethodGet,
	}
}

type MockStore struct {
	invalid []string
}

func (*MockStore) Store(data store.Data) {}
func (*MockStore) GetAll() []store.Data {
	data := make([]store.Data, 0)
	data = append(data, MockData{})

	return data
}
func (s *MockStore) Invalidate(id string) {
	if s.invalid == nil {
		s.invalid = make([]string, 0)
	}
	fmt.Printf("Removing id %s from data store\n", id)

	s.invalid = append(s.invalid, id)
}

func TestExportManager(t *testing.T) {
	ds := &MockStore{}
	ex := &Manager{
		Interval:  5 * time.Second,
		BaseUrl:   "http://localhost:8080",
		DataStore: ds,
	}

	ex.Start(t.Context())

	select {
	case <-time.After(8 * time.Second):
	}

	if !slices.Contains(ds.invalid, metaUuid) {
		t.Fatalf("metaUuid not invalidated!\nmetaUuid: %s\ninvalid: %v\n", metaUuid, ds.invalid)
	}
}
