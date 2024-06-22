package restapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type TRestAPI struct {
	log            *slog.Logger
	server         *http.Server
	enableMwLogger bool
	Router         *chi.Mux
	Generic        Handlers
}

// Creates new server. Set mw to get middleware logger enabled
func New(log *slog.Logger, mw bool) *TRestAPI {
	api := &TRestAPI{log: log, enableMwLogger: mw, Router: chi.NewRouter()}
	api.Router.Use(middleware.RequestID)
	if mw {
		mwLogger := &TLogger{}
		api.Router.Use(mwLogger.New(api.log))
	}
	api.Router.Use(middleware.Recoverer)
	return api
}

// Starts the server (non-blocking)
func (api *TRestAPI) StartAt(addr string) {

	api.server = &http.Server{Addr: addr, Handler: api.Router}

	go func() {
		api.log.Info("starting server",
			slog.Group("settings",
				slog.String("address", addr),
				slog.Bool("mwlogger", api.enableMwLogger),
			),
		)
		if err := api.server.ListenAndServe(); err != nil {
			api.log.Error("failed to start server", "error", err)
		}
	}()

}

// Graceful shutdown with timeout
func (api *TRestAPI) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil {
		api.log.Error("failed to stop server")
		return err
	}

	api.log.Info("server stopped")
	return nil
}
