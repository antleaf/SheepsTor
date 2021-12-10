package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v40/github"
	"net/http"
	"os"
)

func ConfigureRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Throttle(10))
	//TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	//r.Handle("/_resources/assets/*", http.FileServer(http.FS(embeddedAssets)))
	r.Get("/", DefaultHandler)
	//r.Post("/comment", CommentPostHandler)
	r.Post("/update", GitHubWebHookHandler)
	//r.Route("/micropub/{websiteID}", func(r chi.Router) {
	//	r.Use(MicroPubAuthorisationMiddleware)
	//	r.Get("/", MicroPubGetHandler)
	//	r.Post("/", MicroPubPostHandler)
	//	r.Post("/media", MicroPubMediaHandler)
	//})
	//r.Post("/webmention-io/{websiteID}", WebMentionIOHookHandler)
	return r
}

func GitHubWebHookHandler(resp http.ResponseWriter, req *http.Request) {
	logger.Debug("Handling GitHUb webhook post....")
	if req.Method == http.MethodGet {
		if systemReadyToAcceptUpdateRequests == true {
			resp.WriteHeader(http.StatusOK)
		} else {
			resp.WriteHeader(http.StatusServiceUnavailable)
		}
	} else if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(os.Getenv(SheepsTorConfig.GitHubWebHookSecretEnvKey)))
		if err != nil {
			logger.Error(err.Error())
			return
		}
		defer req.Body.Close()
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		switch e := event.(type) {
		case *github.PushEvent:
			websitePtr := SheepsTorConfig.getWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
			if websitePtr != nil {
				ProcessWebsite(websitePtr)
			}
		default:
			logger.Error(fmt.Sprintf("Not a push event %s", github.WebHookType(req)))
			return
		}
	}
}

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello world"))
	resp.WriteHeader(http.StatusOK)
}
