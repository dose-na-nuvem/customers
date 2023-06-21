package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/dose-na-nuvem/customers/proto/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

	lis, freeport, err := ReservaPorta(5, 50051)
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
	// TODO: usar a porta dinamica do socket
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
