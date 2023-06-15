package server

import (
	"context"
	"fmt"

	"github.com/dose-na-nuvem/customers/proto/customer"
	"go.uber.org/zap"
)

// TODO: criar um construtor
type GRPC struct {
	customer.UnimplementedCustomerServer
	logger *zap.Logger
	store  CustomerStore
}

func (g *GRPC) Create(_ context.Context, req *customer.CreateRequest) (*customer.Empty, error) {
	_, err := g.store.CreateCustomer(req.Name)
	if err != nil {
		g.logger.Warn("falha ao criar um customer", zap.Error(err))
		return nil, fmt.Errorf("falha ao criar um customer: %w", err)
	}

	return &customer.Empty{}, nil
}