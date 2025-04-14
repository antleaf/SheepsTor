package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var Debug bool
var log *zap.SugaredLogger

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "--debug=true|false")
}

var rootCmd = &cobra.Command{
	Use: "sheepstor",
	Run: func(cmd *cobra.Command, args []string) {
		initialiseApplication()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
