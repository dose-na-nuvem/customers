package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/dose-na-nuvem/customers/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicializa o servidor do microsserviço.",
	Run: func(cmd *cobra.Command, args []string) {
		h := server.NewHTTP(cfg)

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			cfg.Logger.Info("finalizando o serviço")
			if err := h.Shutdown(context.Background()); err != nil {
				cfg.Logger.Error("erro ao finalizar o serviço", zap.Error(err))
			}
		}()

		cfg.Logger.Info("inicializando o serviço", zap.String("endpoint", cfg.Server.HTTP.Endpoint))
		if err := h.Start(context.Background()); err != nil {
			cfg.Logger.Error("falha ao iniciar o serviço", zap.Error(err))
		}
	},
}
