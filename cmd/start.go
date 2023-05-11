package cmd

import (
	"fmt"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
		cfg.Logger.Info("banco de dados obtido", zap.String("db.type", cfg.DbType))
		server.Serve(cfg)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&cfg.DbType, "dbtype", "<a definir>", "Tipo do banco de dados a ser utilizado.")
}

// Associa cada 'cobra' flag com sua configuração 'viper' (arquivo de configuração e variáveis de ambiente)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determina a convenção de nomes para as flags quando definidas no arquivo de configuração
		configName := f.Name
		// Se usar camelCase no arquivo de configuração, troca hifens com uma string camelCased
		// Como o viper faz comparações case-insensitive (ignorando minúsculas e maiúsculas), não precisamos consertar, só remover os hifens.
		// if replaceHyphenWithCamelCase {
		// 	configName = strings.ReplaceAll(f.Name, "-", "")
		// }

		// Aplica o valor da configuração viper para a flag quando a flag não estiver definida e a do viper tem um valor
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
