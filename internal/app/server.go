package app

import (
	"fmt"
	"net/http"
)

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
