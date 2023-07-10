package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/dose-na-nuvem/customers/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

var errNoTLSConfig = errors.New("servidor sem configuração de TLS")

type HTTP struct {
	logger *zap.Logger

	shutdownCh chan struct{}
	srv        *http.Server

	certFile    string
	certKeyFile string
}

func NewHTTP(cfg *config.Cfg, customerHandler http.Handler) (*HTTP, error) {
	mux := http.NewServeMux()
	mux.Handle("/", otelhttp.NewHandler(customerHandler, "GET /"))

	srv := &http.Server{
		Addr:              cfg.Server.HTTP.Endpoint,
		Handler:           mux,
		ReadHeaderTimeout: cfg.Server.HTTP.ReadHeaderTimeout,
	}

	return HTTPWithServer(cfg, srv)
}

func HTTPWithServer(cfg *config.Cfg, srv *http.Server) (*HTTP, error) {
	if cfg.Server.TLS.CertFile != "" && cfg.Server.TLS.CertKeyFile != "" {
		cfg.Logger.Info("Servidor configurado com opções TLS",
			zap.String("cert_file", cfg.Server.TLS.CertFile),
			zap.String("cert_key_file", cfg.Server.TLS.CertKeyFile),
		)
	} else {
		if cfg.Server.TLS.Insecure {
			cfg.Logger.Info("Servidor sem configurações de TLS! Este servidor está inseguro!")
		} else {
			return nil, errNoTLSConfig
		}
	}

	return &HTTP{
		logger:      cfg.Logger,
		shutdownCh:  make(chan struct{}),
		srv:         srv,
		certFile:    cfg.Server.TLS.CertFile,
		certKeyFile: cfg.Server.TLS.CertKeyFile,
	}, nil
}

func (h *HTTP) Start(_ context.Context) error {
	var err error
	if h.certFile != "" && h.certKeyFile != "" {
		err = h.srv.ListenAndServeTLS(h.certFile, h.certKeyFile)
	} else {
		err = h.srv.ListenAndServe()
	}

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
