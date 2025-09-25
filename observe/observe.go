package observe

import (
	"context"
	"time"

	"github.com/tobiasthedanish/observability-go-sdk/internals/export"
	"github.com/tobiasthedanish/observability-go-sdk/internals/store"
)

type Sdk interface {
	Start() context.CancelFunc
	StartContext(context.Context) context.CancelFunc
}

type sdk struct {
	ex *export.Manager
}

type Config struct {
	/*
		Time between exports. Defaults to 30 seconds.
	*/
	ExportInterval time.Duration
	/*
	 Base url of your observability server. Including protcol like "https://...".
	 Defaults to "http://localhost:8080"
	*/
	ExportBaseUrl string
}

func InitSdk(config Config) Sdk {
	config = setDefaultConfig(config)

	ex := &export.Manager{
		Interval:  config.ExportInterval,
		BaseUrl:   config.ExportBaseUrl,
		DataStore: &store.DummyStore{},
	}

	return &sdk{
		ex: ex,
	}
}

func setDefaultConfig(config Config) Config {
	if config.ExportInterval == 0 {
		config.ExportInterval = 30 * time.Second
	}
	if config.ExportBaseUrl == "" {
		config.ExportBaseUrl = "http://localhost:8080"
	}

	return config
}

func (o *sdk) Start() context.CancelFunc {
	return o.StartContext(context.Background())
}

func (o *sdk) StartContext(c context.Context) context.CancelFunc {
	ctx, cancel := context.WithCancel(c)

	o.ex.Start(ctx)

	return cancel
}
