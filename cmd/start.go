package cmd

import (
	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	cfg = config.New()
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicializa o servidor do microsserviço.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Logger.Info("banco de dados obtido", zap.String("db.type", cfg.DBType))
		server.Serve(cfg)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&cfg.DBType, "dbtype", "<a definir>", "Tipo do banco de dados a ser utilizado.")
}
