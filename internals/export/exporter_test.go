package export

import (
	"net/http"
	"testing"
	"time"
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

	jobChan <- ExportJob{
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
	case <-time.After(timeout):
		t.Fatalf("Context timeout before response\n")
	}
}
