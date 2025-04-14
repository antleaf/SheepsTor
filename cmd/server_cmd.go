package cmd

import (
	"fmt"
	. "github.com/antleaf/sheepstor/pkg"
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
	Router = NewRouter()
	Renderer = NewRenderer()
	log.Infof("Running as HTTP Process on port %d", Config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", Config.Port), Router)
	if err != nil {
		log.Error(err.Error())
	}
}
