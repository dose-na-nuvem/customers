package server //nolint

import (
	"testing"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/stretchr/testify/assert"
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
	srv, err := NewHTTP(cfg)

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
				assert.Contains(t, "Servidor sem configurações de TLS! Este servidor está inseguro!", logs.All()[0].Message)
			} else {
				assert.Equal(t, errNoTLSConfig, err)
			}
		})
	}
}
