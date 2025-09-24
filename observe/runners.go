package observe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tobiasthedanish/observability-go-sdk/internals/export"
)

func RunExportTest() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobChan := make(chan export.ExportJob)
	resChan := make(chan export.ExportResult)

	ex := &export.Exporter{
		Jobs: jobChan,
		Res:  resChan,
	}

	ex.Start(ctx)

	jobChan <- export.ExportJob{
		Method:  http.MethodGet,
		Url:     "http://localhost:8080/health",
		Headers: http.Header{},
		// data: struct{ Message string }{
		// 	Message: "This is a new job",
		// },
	}

	select {
	case res := <-resChan:
		fmt.Printf("Received result: %+v\n", res)
	case <-time.After(10 * time.Second):
		fmt.Printf("Timeout without res\n")
	}
}

func RunSdk() {
	sdk := InitSdk(Config{
		ExportInterval: 2 * time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())
	sdk.StartContext(ctx)

	select {
	case <-time.After(3 * time.Second):
		cancel()
	}
}
