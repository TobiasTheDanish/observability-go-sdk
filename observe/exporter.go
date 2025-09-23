package observe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type exportJob struct {
	method  string
	url     string
	headers http.Header
	data    any
}

type exporter struct {
	isRunning bool
	// channel to receive jobs for export
	jobs chan exportJob
}

func (e *exporter) start(ctx context.Context) {
	if e.jobs == nil {
		e.jobs = make(chan exportJob)
	}
	if e.isRunning {
		// already started... returning early
		return
	}

	e.isRunning = true
	go func() {
		for {
			select {
			case job := <-e.jobs:
				handleJob(job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func handleJob(job exportJob) {
	fmt.Printf("Received job: %+v\n", job)
	jsonData, err := json.Marshal(job.data)
	if err != nil {
		fmt.Printf("Marshalling data as json failed: %v\n", err)
		return
	}
	fmt.Printf("JSON data: %v\n", string(jsonData))

	req, err := http.NewRequest(job.method, job.url, bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Creating new http request failed: %v\n", err)
		return
	}
	fmt.Printf("request has been created\n")

	for key, values := range job.headers {
		for _, val := range values {
			req.Header.Add(key, val)
		}
	}
	fmt.Printf("headers have been set\n")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Sending http request failed: %v\n", err)
		return
	}
	fmt.Printf("request have been sent. Response status: %s\n", res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reading http response body failed: %v\n", err)
		return
	}

	fmt.Printf("Response body: %v\n", body)
}
