package cmd

import (
	"fmt"
	. "github.com/antleaf/SheepsTor/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"net/http"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	//updateWebsiteCmd.Flags().StringVarP(&sites, "sites", "", "", "--sites all|<some_id>")
}

var router chi.Router

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		initialiseApplication()
		runServer()
	},
}

func runServer() {
	router = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(10))
	//TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	router.Get("/", DefaultHandler)
	router.Post("/update", registry.GitHubWebHookHandler)
	Log.Infof("Running as HTTP Process on port %d", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
	if err != nil {
		Log.Error(err.Error())
	}
}

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello world"))
	resp.WriteHeader(http.StatusOK)
}
