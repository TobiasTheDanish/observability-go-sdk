package export

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ExportJob struct {
	Method  string
	Url     string
	Headers http.Header
	Data    any
}

type ExportResult struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type Exporter struct {
	isRunning bool
	// channel to receive Jobs for export
	Jobs <-chan ExportJob
	Res  chan<- ExportResult
}

func (e *Exporter) Start(ctx context.Context) {
	if e.isRunning {
		// already started... returning early
		return
	}
	if e.Jobs == nil {
		e.Jobs = make(chan ExportJob)
	}
	if e.Res == nil {
		e.Res = make(chan ExportResult)
	}

	e.isRunning = true
	go e.run(ctx)
}

func (e *Exporter) run(ctx context.Context) {
	for {
		select {
		case job := <-e.Jobs:
			go func() {
				res, ok := handleJob(job, ctx)
				if ok {
					e.Res <- res
					fmt.Printf("ExportResult sent to result channel: %+v\n", res)
				}
			}()
		case <-ctx.Done():
			fmt.Printf("Exporter context has been canceled. Stopping\n")
			return
		}
	}
}

func handleJob(job ExportJob, ctx context.Context) (res ExportResult, ok bool) {
	var jsonData []byte
	var err error

	if job.Data != nil {
		jsonData, err = json.Marshal(job.Data)
		if err != nil {
			fmt.Printf("Marshalling data as json failed: %v\n", err)
			return
		}
	}

	req, err := http.NewRequestWithContext(ctx, job.Method, job.Url, bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Creating new http request failed: %v\n", err)
		return
	}

	for key, values := range job.Headers {
		for _, val := range values {
			req.Header.Add(key, val)
		}
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Sending http request failed: %v\n", err)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Reading http response body failed: %v\n", err)
		return
	}

	ok = true
	res = ExportResult{
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       body,
	}
	return
}
