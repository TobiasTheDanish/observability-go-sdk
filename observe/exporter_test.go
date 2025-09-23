package observe

import (
	"net/http"
	"testing"
)

func TestExporter(t *testing.T) {
	ctx := t.Context()
	jobChan := make(chan exportJob)

	ex := &exporter{
		jobs: jobChan,
	}

	ex.start(ctx)

	jobChan <- exportJob{
		url:     "http://localhost:8080/hello",
		headers: http.Header{},
		data: struct{ Message string }{
			Message: "This is a new job",
		},
	}
}
