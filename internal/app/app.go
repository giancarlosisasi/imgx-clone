package app

import "github.com/giancarlosisasi/imgix-clone/internal/config"

type Application struct {
	config *config.Config
}

func NewApp(config *config.Config) *Application {
	return &Application{
		config: config,
	}
}
