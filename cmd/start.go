package cmd

import (
	"fmt"
	"strings"

	"github.com/dose-na-nuvem/customers/config"
	"github.com/dose-na-nuvem/customers/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfg = config.New()

	defaultConfigFilename = "config"
	defaultConfigPath     = "."
	envPrefix             = "DOSE"
	// Troca flags que tenham hifen com camelCase no arquivo de configuração.
	replaceHyphenWithCamelCase = true
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicializa o servidor do microsserviço.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// É possível associar cobra e viper em outros lugares, mas PersistencePreRunE na raiz do comando funciona muito bem
		return configViperCfgWithCobraFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// os flags compostos devem usar a forma com sublinhado
		cfg.Logger.Info("banco de dados obtido", zap.String("db_type", cfg.DbType))
		server.Serve(cfg)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&cfg.DbType, "db_type", "<a definir>", "Tipo do banco de dados a ser utilizado.")
}

// Configura um comando do cobra com ajustes necessários ao sincronismo com o viper.
func configViperCfgWithCobraFlags(cmd *cobra.Command) error {
	v := viper.New()

	// Configura o nome padrão do arquivo de configuração, sem a extensão.
	v.SetConfigName(defaultConfigFilename)

	// Configure quantos caminhos forem necessários para o viper buscar o arquivo
	// Nesse caso especificamente, vamos considerar somente o diretório de trabalho.
	v.AddConfigPath(defaultConfigPath)

	// Tenta ler o arquivo de configuração, ignorando erros caso o mesmo não seja encontrado
	// Retorna um erro se não conseguirmos analisar o arquivo de configuração encontrado.
	if err := v.ReadInConfig(); err != nil {
		// Não há problemas se não existir um arquivo de configuração.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Variáveis de ambiente não podem ter hifens, então associaremos a uma variável equivalente usando sublinhado
	// e.g. --cor-favorita será COR_FAVORITA
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Quando associamos flags às variáveis de ambiente, esperamos que as mesmas estejam prefixadas de alguma forma
	// e.g. uma flag como --numero associada a uma variavel de ambiente DOSE_NUMERO
	// Isso evita muitos conflitos quando se executam diversos aplicativos num mesmo ambiente
	v.SetEnvPrefix(envPrefix)

	// Confirma o uso de variáveis de ambiente
	// Funciona muito bem para nomes simples, mas precisa de ajustes para nomes compostos como --cor-favorita
	// que serão ajustados na função bindFlags
	v.AutomaticEnv()

	// Associa as flags do comando ao viper
	bindFlags(cmd, v)

	return nil
}

// Associa cada 'cobra' flag com sua configuração 'viper' (arquivo de configuração e variáveis de ambiente).
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determina a convenção de nomes para as flags quando definidas no arquivo de configuração
		configName := f.Name
		// Se usar camelCase no arquivo de configuração, troca hifens com uma string camelCased
		// Como o viper faz comparações case-insensitive (ignorando minúsculas e maiúsculas), não precisamos consertar, só remover os hifens.
		if replaceHyphenWithCamelCase {
			configName = strings.ReplaceAll(f.Name, "-", "")
		}

		// Aplica o valor da configuração viper para a flag quando a flag não estiver definida e a do viper tem um valor
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
