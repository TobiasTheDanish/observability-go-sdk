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
	BaseUrl   string
	DataStore store.DataStore[store.Data]

	exporter *Exporter
	jobChan  chan ExportJob
	resChan  chan ExportResult
}

func (m *Manager) Start(ctx context.Context) {
	if m.jobChan == nil {
		m.jobChan = make(chan ExportJob)
	}
	if m.resChan == nil {
		m.resChan = make(chan ExportResult)
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
	m.startResultHandling(ctx)
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
		meta := data.Metadata()

		m.jobChan <- ExportJob{
			DataId:  meta.Id,
			Method:  meta.ExportMethod,
			Url:     fmt.Sprintf("%s%s", m.BaseUrl, meta.ExportPath),
			Headers: make(http.Header),
			Data:    data,
		}
	}
}

func (m *Manager) startResultHandling(ctx context.Context) {
	go func() {
		for {
			select {
			case res := <-m.resChan:
				m.handleResult(res)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *Manager) handleResult(res ExportResult) {
	fmt.Printf("Received result: %+v\n", res)

	switch {
	case 200 <= res.StatusCode && res.StatusCode <= 299 ||
		400 <= res.StatusCode && res.StatusCode < 429:
		m.DataStore.Invalidate(res.DataID)
	}
}
