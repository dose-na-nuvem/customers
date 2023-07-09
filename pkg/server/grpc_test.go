package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/dose-na-nuvem/customers/proto/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreateCustomer(t *testing.T) {
	// prepare
	var called bool
	st := &mockStore{createCustomerFunc: func(name string) (*model.Customer, error) {
		called = true

		assert.Equal(t, "Fulano de Tal", name)

		return nil, nil
	}}

	lis, freeport, err := GetListenerWithFallback(5, 50051)
	require.NoError(t, err)
	defer lis.Close()

	// TODO: ver quais ServerOption podemos colocar como propriedades no arquivo de configuração
	s := grpc.NewServer()
	defer s.Stop()

	customer.RegisterCustomerServer(s, &GRPC{
		logger: zap.NewNop(),
		store:  st,
	})

	// TODO: achar um jeito de fazer isso sem bloquear
	go func() {
		err = s.Serve(lis)
		require.NoError(t, err)
	}()

	// prepara a parte de cliente
	endpoint := fmt.Sprintf("localhost:%d", freeport)
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	cl := customer.NewCustomerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := &customer.CreateRequest{
		Name: "Fulano de Tal",
	}
	// test
	resp, err := cl.Create(ctx, c)

	// verify
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, called)
}

func TestGRPCServerTLS(t *testing.T) {
	// prepare
	core, _ := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	testCases := []struct {
		desc        string
		cfg         *config.Cfg
		shouldErr   bool
		certFile    string
		certKeyFile string
	}{
		{
			desc:        "has certs, insecure is set to false",
			cfg:         config.New(),
			shouldErr:   false,
			certFile:    "fixtures/certs/cert.pem",
			certKeyFile: "fixtures/certs/cert-key.pem",
		},
		{
			desc:        "has broken certs",
			cfg:         config.New(),
			shouldErr:   true,
			certFile:    "fixtures/certs/invalid.pem",
			certKeyFile: "fixtures/certs/invalid-key.pem",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// prepare
			cfg := tC.cfg
			cfg.Logger = logger
			if tC.certFile != "" {
				cfg.Server.TLS.CertFile = tC.certFile
				cfg.Server.TLS.CertKeyFile = tC.certKeyFile
			}

			// test
			computed, err := buildServerOptions(tC.cfg)

			// assert
			if !tC.shouldErr {
				require.NoError(t, err)
				assert.NotEmpty(t, computed, "esperava-se configurações preenchidas")
			} else {
				assert.Empty(t, computed, "não se espera configuração alguma")
				assert.Error(t, err)
			}
		})
	}
}

func TestGRPCServerInsecure(t *testing.T) {
	// prepare
	core, _ := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	testCases := []struct {
		desc        string
		cfg         *config.Cfg
		setInsecure bool
		insecure    bool
		shouldErr   bool
	}{
		{
			desc:        "1 no certs, insecure not set",
			cfg:         config.New(),
			setInsecure: false,
			insecure:    false,
			shouldErr:   true,
		},
		{
			desc:        "2 no certs, insecure is set to true",
			cfg:         config.New(),
			setInsecure: true,
			insecure:    true,
			shouldErr:   false,
		},
		{
			desc:        "3 no certs, insecure is set to false",
			cfg:         config.New(),
			setInsecure: true,
			insecure:    false,
			shouldErr:   true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// prepare
			cfg := tC.cfg
			cfg.Logger = logger
			if tC.setInsecure {
				cfg.Server.TLS.Insecure = tC.insecure
			}

			// test
			computed, err := buildServerOptions(tC.cfg)

			// assert
			if !tC.shouldErr {
				require.NoError(t, err)
				if tC.insecure {
					assert.Empty(t, computed, "esperava configuração vazia")
				} else {
					assert.NotEmpty(t, computed, "esperava-se configurações preenchidas")
				}
			} else {
				assert.Empty(t, computed, "não se espera configuração alguma")
				assert.Error(t, err)
			}
		})
	}
}

func TestGRPC_NonBlockingStartSuccessful(t *testing.T) {
	// arrange
	ctx := context.Background()
	errorChannel := make(chan error, 1) // buffered

	lis, _, err := GetListenerWithFallback(3, 40404)
	require.NoError(t, err, "não foi possivel usar uma porta livre ")

	g := &GRPC{
		logger:   zap.NewNop(),
		store:    &mockStore{},
		grpc:     grpc.NewServer(),
		listener: lis,
	}

	// act
	g.Start(ctx, errorChannel)

	// assert
	assert.Empty(t, errorChannel, "o grpc iniciou com sucesso")
	assert.Eventually(t, func() bool {
		return assert.Empty(t, errorChannel, "o grpc iniciou com sucesso")
	}, 300*time.Millisecond, 20*time.Millisecond, "o http falhou")

	// assert
	// time.Sleep(time.Second * 1)
	err = g.Shutdown(ctx)
	assert.NoError(t, err, "não deve ter erro se foi inicializado corretamente")
}
