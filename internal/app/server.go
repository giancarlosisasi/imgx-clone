package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func (app *Application) Serve() error {
	defer app.ctxCancel()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Port),
		Handler: app.routes(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		// ctrl + c
		syscall.SIGINT,
		// kubernetes/docker stop
		syscall.SIGTERM,
	)

	// channel to receive server errors
	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Msgf("> Server running in: http://localhost:%d", app.config.Port)

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		// catch any server error (which means the server never started)
		return err
	case <-quit:
		// shutdown signal received, continue to graceful shutdown
	}

	log.Warn().Msg("⚠️  Shutdown signal received, starting graceful shutdown...")

	// cancel context to signal all background jobs to stop
	app.ctxCancel()

	// create a timeout context for shutdown, only wait for 30seconds
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// shutdown http server (stops accepting new connections)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
		return err
	} else {
		log.Info().Msg("✅ HTTP server stopped accepting new requests")
	}

	// wait for all background jobs to finish
	log.Info().Msg("⏳ Waiting for background jobs to complete...")

	done := make(chan struct{})
	go func() {
		app.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("✅ All background jobs completed")
	case <-shutdownCtx.Done():
		log.Warn().Msg("⚠️  Timeout reached, forcing shutdown")
	}

	log.Info().Msg("👋 Server shutdown complete")

	return nil
}
