package main

import (
	"github.com/giancarlosisasi/imgix-clone/internal/app"
	"github.com/giancarlosisasi/imgix-clone/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	config := config.NewConfig()
	app := app.NewApp(config)

	err := app.Serve()
	if err != nil {
		log.Fatal().Err(err).Msg("There was an error when starting the app")
	}

}
