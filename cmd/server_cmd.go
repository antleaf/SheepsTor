package cmd

import (
	"fmt"
	. "github.com/antleaf/SheepsTor/internal"
	"github.com/spf13/cobra"
	"net/http"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		initialiseApplication()
		runServer()
	},
}

func runServer() {
	InitialiseServer()
	Log.Infof("Running as HTTP Process on port %d", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), Router)
	if err != nil {
		Log.Error(err.Error())
	}
}
