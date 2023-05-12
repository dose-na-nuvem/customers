package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&cfg.DBType, "dbtype", "<a definir>", "Tipo do banco de dados a ser utilizado.")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
