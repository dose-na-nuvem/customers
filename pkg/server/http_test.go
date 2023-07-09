package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewHTTP(t *testing.T) {
	// prepare
	cfg := config.New()
	cfg.Server = config.ServerSettings{
		TLS: &config.TLSSettings{
			CertFile:    "fixtures/cert.pem",
			CertKeyFile: "fixtures/cert-key.pem",
		},
	}

	// test
	srv, err := NewHTTP(cfg, NewCustomerHandler(cfg.Logger, nil))

	// verify
	assert.NoError(t, err)
	assert.Equal(t, "fixtures/cert.pem", srv.certFile)
	assert.Equal(t, "fixtures/cert-key.pem", srv.certKeyFile)
}

func TestHTTPWithInsecureServer(t *testing.T) {
	// prepare
	testCases := []struct {
		desc     string
		insecure bool
	}{
		{
			desc:     "no cert, insecure not set",
			insecure: false,
		},
		{
			desc:     "no cert, insecure is set",
			insecure: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// prepare
			core, logs := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			cfg := config.New()
			cfg.Logger = logger
			cfg.Server.TLS.Insecure = tC.insecure

			// test
			_, err := HTTPWithServer(cfg, nil)

			// verify
			if tC.insecure {
				assert.Len(t, logs.All(), 1)
				assert.Contains(t, "Servidor HTTP sem configurações de TLS! Este servidor está inseguro!", logs.All()[0].Message)
			} else {
				assert.Equal(t, errNoTLSConfig, err)
			}
		})
	}
}

// flaky test
func TestHTTP_NonBlockingStartSuccessful(t *testing.T) {
	var err error
	// prepare
	ctx := context.Background()
	errChannel := make(chan error)

	listener, port, err := GetListenerWithFallback(3, 43678)
	require.NoError(t, err, "não foi possivel alocar uma porta livre")
	listener.Close()
	freePortEndpoint := fmt.Sprintf("localhost:%d", port)

	cfg := config.New()
	cfg.Server.HTTP.ReadHeaderTimeout = 1 * time.Second
	cfg.Server.HTTP.Endpoint = freePortEndpoint

	srv := &http.Server{
		Addr:              freePortEndpoint,
		ReadHeaderTimeout: 2 * time.Second,
	}

	h := &HTTP{
		logger:      cfg.Logger,
		srv:         srv,
		certFile:    cfg.Server.TLS.CertFile,
		certKeyFile: cfg.Server.TLS.CertKeyFile,
	}

	// act
	h.Start(ctx, errChannel)

	// assert
	assert.Eventually(t, func() bool {
		return assert.Empty(t, errChannel, "o http iniciou com sucesso")
	}, 300*time.Millisecond, 20*time.Millisecond, "o http falhou")

	// assert
	err = h.Shutdown(ctx)
	assert.NoError(t, err, "não deve ter erro se foi inicializado corretamente")
}
