package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var port int

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "", 8081, "--port=8081")
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
	log.Infof("Running as HTTP Process on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), Router)
	if err != nil {
		log.Error(err.Error())
	}
}
