package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/dose-na-nuvem/customers/pkg/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const MaxServers = 2

type Customer struct {
	cfg               *config.Cfg
	srv               *server.HTTP
	grpc              *server.GRPC
	asyncErrorChannel chan error
	signalsChannel    chan os.Signal
}

func New(cfg *config.Cfg) *Customer {
	return &Customer{
		cfg:               cfg,
		asyncErrorChannel: make(chan error, MaxServers), // buffered
		signalsChannel:    make(chan os.Signal),
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
	go func() {
		if grpcErr := c.grpc.Start(ctx); !errors.Is(grpcErr, grpc.ErrServerStopped) {
			c.cfg.Logger.Error("falha ao iniciar o servidor GRPC", zap.Error(grpcErr))
			c.asyncErrorChannel <- grpcErr
		}
	}()

	c.srv, err = server.NewHTTP(c.cfg, ch)
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor HTTP: %w", err)
	}

	go func() {
		if httpErr := c.srv.Start(ctx); !errors.Is(httpErr, http.ErrServerClosed) {
			c.cfg.Logger.Error("falha ao iniciar o servidor HTTP", zap.Error(httpErr))
			c.asyncErrorChannel <- httpErr
		}
	}()

	signal.Notify(c.signalsChannel, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(c.signalsChannel)

LOOP:
	for {
		select {
		case err := <-c.asyncErrorChannel:
			c.cfg.Logger.Error("falha ao iniciar o servidor: %w", zap.Error(err))
			break LOOP
		case signal := <-c.signalsChannel:
			c.cfg.Logger.Debug("signal received", zap.Any("signal", signal.String()))
			err := c.Shutdown(ctx)
			if err != nil {
				c.cfg.Logger.Error("falha ao finalizar servidor: %w", zap.Error(err))
			}
			return nil
		}
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
