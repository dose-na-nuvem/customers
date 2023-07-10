package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/dose-na-nuvem/customers/pkg/server"
	"github.com/dose-na-nuvem/customers/pkg/telemetry"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
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

func (c *Customer) bootstrap(ctx context.Context) (server.CustomerStore, error) {
	// este rastreador é somente para o processo de bootstrapping

	tr := otel.GetTracerProvider().Tracer("bootstrap")

	ctx, rootSpan := tr.Start(ctx, "bootstrap")

	_, span := tr.Start(ctx, "db/open")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}
	span.End()

	// Migrate the schema
	_, span = tr.Start(ctx, "db/migrate")
	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		return nil, fmt.Errorf("falha ao migrar o esquema do banco de dados: %w", err)
	}
	span.End()

	rootSpan.End()

	return model.NewStore(db), nil
}

func (c *Customer) Start(ctx context.Context) error {
	var err error

	tp, err := telemetry.NewTracerProvider()
	if err != nil {
		return fmt.Errorf("falha ao iniciar os rastreadores: %w", err)
	}
	otel.SetTracerProvider(tp)

	store, err := c.bootstrap(ctx)
	if err != nil {
		return err
	}
	ch := server.NewCustomerHandler(c.cfg.Logger, store)

	c.grpc, err = server.NewGRPC(c.cfg, store)
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor GRPC: %w", err)
	}
	c.grpc.Start(ctx, c.asyncErrorChannel)

	c.srv, err = server.NewHTTP(c.cfg, ch)
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor HTTP: %w", err)
	}
	c.srv.Start(ctx, c.asyncErrorChannel)

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

	if err := c.grpc.Shutdown(ctx); err != nil {
		return fmt.Errorf("erro ao finalizar o serviço: %w", err)
	}

	return nil
}
