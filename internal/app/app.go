package app

import (
	"context"
	"sync"

	"github.com/giancarlosisasi/imgix-clone/internal/config"
)

type Application struct {
	config    *config.Config
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewApp(config *config.Config) *Application {
	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		config:    config,
		wg:        sync.WaitGroup{},
		ctx:       ctx,
		ctxCancel: cancel,
	}
}
