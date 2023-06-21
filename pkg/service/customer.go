package service

import (
	"context"
	"fmt"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/dose-na-nuvem/customers/pkg/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Customer struct {
	cfg  *config.Cfg
	srv  *server.HTTP
	grpc *server.GRPC
}

func New(cfg *config.Cfg) *Customer {
	return &Customer{
		cfg: cfg,
	}
}

func (c *Customer) Start(ctx context.Context) error {
	var err error

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		return fmt.Errorf("falha ao migrar o esquema do banco de dados: %w", err)
	}

	store := model.NewStore(db)

	ch := server.NewCustomerHandler(c.cfg.Logger, store)

	c.grpc, err = server.NewGRPC(c.cfg, store)
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor GRPC: %w", err)
	}
	if err = c.grpc.Start(ctx); err != nil {
		return fmt.Errorf("falha ao inicia o servidor GRPC: %w", err)
	}

	c.srv, err = server.NewHTTP(c.cfg, ch)
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor HTTP: %w", err)
	}

	if err = c.srv.Start(ctx); err != nil {
		return fmt.Errorf("falha ao iniciar o servidor HTTP: %w", err)
	}

	return nil
}

func (c *Customer) Shutdown(ctx context.Context) error {
	if err := c.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("erro ao finalizar o serviço: %w", err)
	}

	// TODO: finalizar o servidor gRPC
	if err := c.grpc.Shutdown(ctx); err != nil {
		return fmt.Errorf("erro ao finalizar o serviço: %w", err)
	}

	return nil
}
