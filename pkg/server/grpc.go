package server

import (
	"context"
	"fmt"
	"net"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/proto/customer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TODO: criar um construtor
type GRPC struct {
	customer.UnimplementedCustomerServer
	logger   *zap.Logger
	store    CustomerStore
	grpc     *grpc.Server
	listener net.Listener
}

func (g *GRPC) Create(_ context.Context, req *customer.CreateRequest) (*customer.Empty, error) {
	_, err := g.store.CreateCustomer(req.Name)
	if err != nil {
		g.logger.Warn("falha ao criar um customer", zap.Error(err))
		return nil, fmt.Errorf("falha ao criar um customer: %w", err)
	}

	return &customer.Empty{}, nil
}

func NewGRPC(cfg *config.Cfg, store CustomerStore) (*GRPC, error) {
	lis, err := net.Listen("tcp", cfg.Server.GRPC.Endpoint)
	if err != nil {
		return &GRPC{}, err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	grpc := &GRPC{
		logger:   cfg.Logger,
		store:    store,
		grpc:     grpcServer,
		listener: lis,
	}

	customer.RegisterCustomerServer(grpcServer, grpc)

	return grpc, nil
}

func (g *GRPC) Start(_ context.Context) error {
	g.logger.Info("iniciando servidor gRPC")
	err := g.grpc.Serve(g.listener)
	return err
}

func (g *GRPC) Shutdown(_ context.Context) error {
	g.logger.Info("finalizando servidor gRPC")
	g.grpc.GracefulStop()
	err := g.listener.Close()
	if err != nil {
		return err
	}
	return nil
}
