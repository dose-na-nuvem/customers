package cmd

import (
	"github.com/dose-na-nuvem/customers/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicializa o servidor do microsserviço.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Logger.Info("banco de dados obtido",
			zap.String("db.type", cfg.Database.Type),
			zap.String("db.username", cfg.Database.Username),
		)
		server.Serve(cfg)
	},
}
