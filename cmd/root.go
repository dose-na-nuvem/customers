package cmd

import (
	"os"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	configFile string
	cfg        = config.New()
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use: "customers",

	//nolint
	Short: "Microsserviço responsável pelo gerenciamento de clientes.",

	//nolint
	Long: `Microsserviço responsável pelo gerenciamento de clientes.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(startCmd)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml",
		"Define o arquivo de configuração a utilizar.")

	startCmd.Flags().StringVar(&cfg.Database.Type, "db.type", "<a definir>",
		"Tipo do banco de dados a ser utilizado.")

	startCmd.Flags().StringVar(&cfg.Database.Username, "db.username", "<a definir>",
		"Usuário do banco de dados a ser utilizado.")

	startCmd.Flags().StringVar(&cfg.Database.Password, "db.password", "<a definir>",
		"Senha do usuário do banco de dados a ser utilizado.")

	startCmd.Flags().StringVar(&cfg.Server.HTTP.Endpoint, "server.http.endpoint", "localhost:56433",
		"Endereço onde o serviço vai servir requisições.")

	startCmd.Flags().Int64Var(&cfg.Server.HTTP.ReadHeaderTimeout, "server.http.readHeaderTimeout", 1000, "Tempo máximo de leitura dos headers de uma requisição em Milissegundos")

	// tie Viper to flags
	if err := viper.BindPFlags(startCmd.Flags()); err != nil {
		cfg.Logger.Error("falha ao ligar as flags", zap.Error(err))
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	// Configura o nome padrão do arquivo de configuração, sem a extensão.
	viper.SetConfigFile(configFile)

	// Tenta ler o arquivo de configuração, ignorando erros caso o mesmo não seja encontrado
	// Retorna um erro se não conseguirmos analisar o arquivo de configuração encontrado.
	if err := viper.ReadInConfig(); err != nil {
		// Não há problems se não existir um arquivo de configuração.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cfg.Logger.Error("arquivo não encontrado",
				zap.String("arquivo", configFile),
				zap.Error(err),
			)

			return
		}

		cfg.Logger.Error("falha na leitura do arquivo de configuração", zap.Error(err))
	} else {
		cfg.Logger.Info("arquivo de configuração lido", zap.String("config", configFile))
	}

	// convert Viper's internal state into our configuration object
	if err := viper.Unmarshal(&cfg); err != nil {
		cfg.Logger.Error("falhou ao converter o arquivo de configuração", zap.Error(err))

		return
	}
}
