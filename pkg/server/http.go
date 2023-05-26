package server

import (
	"context"
	"net/http"
	"time"

	"github.com/dose-na-nuvem/customers/config"
	"go.uber.org/zap"
)

type HTTP struct {
	logger *zap.Logger

	shutdownCh chan struct{}
	srv        *http.Server
}

func NewHTTP(cfg *config.Cfg) *HTTP {
	mux := http.NewServeMux()
	mux.Handle("/", NewCustomerHandler(cfg.Logger))

	srv := &http.Server{
		Addr:              cfg.Server.HTTP.Endpoint,
		Handler:           mux,
		ReadHeaderTimeout: cfg.Server.HTTP.ReadHeaderTimeout,
	}

	return HTTPWithServer(cfg, srv)
}

func HTTPWithServer(cfg *config.Cfg, srv *http.Server) *HTTP {
	return &HTTP{
		logger:     cfg.Logger,
		shutdownCh: make(chan struct{}),
		srv:        srv,
	}
}

func (h *HTTP) Start(_ context.Context) error {
	err := h.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (h *HTTP) Shutdown(ctx context.Context) error {
	// We received an interrupt signal, shut down.
	if err := h.srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		return err
	}

	close(h.shutdownCh)

	return nil
}
