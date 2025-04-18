package cmd

import (
	. "github.com/antleaf/sheepstor/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/render"
	"net/http"
)

var Router chi.Router
var Renderer *render.Render

func NewRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(10))
	//TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	router.Get("/", DefaultHandler)
	router.Post("/update", GitHubWebHookHandler)
	router.Handle("/assets/*", http.FileServer(http.FS(embeddedAssets)))
	return router
}
