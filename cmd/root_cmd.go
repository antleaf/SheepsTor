package cmd

import (
	. "github.com/antleaf/SheepsTor/internal"
	"github.com/spf13/cobra"
	"os"
)

var config = Configuration{}
var Debug bool
var registry WebsiteRegistry

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
