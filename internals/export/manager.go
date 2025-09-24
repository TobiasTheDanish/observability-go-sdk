package export

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tobiasthedanish/observability-go-sdk/internals/store"
)

type Manager struct {
	Interval  time.Duration
	DataStore store.DataStore[any]

	exporter *Exporter
	jobChan  chan ExportJob
	resChan  chan ExportResult
}

func (m *Manager) Start(ctx context.Context) {
	if m.jobChan == nil {
		m.jobChan = make(chan ExportJob)
	}
	if m.resChan == nil {
		m.resChan = make(chan ExportResult, 1)
	}
	if m.exporter == nil {
		m.exporter = &Exporter{
			Jobs: m.jobChan,
			Res:  m.resChan,
		}
	}
	if !m.exporter.isRunning {
		m.exporter.Start(ctx)
	}

	m.startInterval(ctx)
}

func (m *Manager) startInterval(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(m.Interval):
				m.handleExport()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *Manager) handleExport() {
	fmt.Printf("Export triggered\n")
	all := m.DataStore.GetAll()
	for _, data := range all {
		m.jobChan <- ExportJob{
			Method:  http.MethodGet,
			Url:     "http://localhost:8080/health",
			Headers: make(http.Header),
			Data:    data,
		}
	}
}
