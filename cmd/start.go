package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/dose-na-nuvem/customers/pkg/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicializa o servidor do microsserviço.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		svc := service.New(cfg)

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			cfg.Logger.Info("finalizando o serviço")

			// TODO: colocar uma deadline para o shutdown
			if err := svc.Shutdown(ctx); err != nil {
				cfg.Logger.Error("erro ao finalizar o serviço: %w", zap.Error(err))
			}
			cfg.Logger.Info("serviço finalizado com sucesso")
		}()

		// TODO: ajustar as msgs conforme o servidor que possa ter dado problema
		cfg.Logger.Info("inicializando o serviço",
			zap.String("http", cfg.Server.HTTP.Endpoint),
			zap.String("grpc", cfg.Server.GRPC.Endpoint))
		if err := svc.Start(ctx); err != nil {
			cfg.Logger.Error("erro ao inicializar o serviço", zap.Error(err))
		}
	},
}
