package export

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExporter(t *testing.T) {
	ctx := t.Context()
	jobChan := make(chan ExportJob)
	resChan := make(chan ExportResult)

	ex := &Exporter{
		Jobs: jobChan,
		Res:  resChan,
	}

	ex.Start(ctx)

	id := uuid.New().String()
	jobChan <- ExportJob{
		DataId:  id,
		Method:  "GET",
		Url:     "http://localhost:8080/health",
		Headers: http.Header{},
		Data: struct{ Message string }{
			Message: "This is a new job",
		},
	}

	deadline, ok := t.Context().Deadline()
	var timeout time.Duration
	if ok {
		timeout = deadline.Sub(time.Now())
	} else {
		timeout = 10 * time.Second
	}
	select {
	case res := <-resChan:
		t.Logf("Received response: %+v\n", res)
		if res.DataID != id {
			t.Fatalf("Received data id in response different from data provided in job.\nres: %s\njob: %s\n", res.DataID, id)
		}
	case <-time.After(timeout):
		t.Fatalf("Context timeout before response\n")
	}
}
