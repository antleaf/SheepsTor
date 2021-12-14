package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v40/github"
	"io/ioutil"
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
	r.Post("/webmentionio/{websiteID}", WebMentionIOHookHandler)
	//r.Route("/micropub/{websiteID}", func(r chi.Router) {
	//	r.Use(MicroPubAuthorisationMiddleware)
	//	r.Get("/", MicroPubGetHandler)
	//	r.Post("/", MicroPubPostHandler)
	//	r.Post("/media", MicroPubMediaHandler)
	//})
	return r
}

func WebMentionIOHookHandler(w http.ResponseWriter, r *http.Request) {
	websiteID := chi.URLParam(r, "websiteID")
	logger.Debugf("Received webmention.io webhook for website with ID: %s", websiteID)
	website := SheepsTorConfig.getWebsiteByID(websiteID)
	if website == nil {
		logger.Error(errors.New("no website found with ID: " + websiteID))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	payloadJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	webmentionPayload := WebmentionIOPayload{}
	webmention, err := webmentionPayload.LoadAndValidate(payloadJson, website.SheepsTorProcessing.WebmentionIoWebhookSecret)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageNode := website.SiteMap.GetNodeByPermalink(webmention.Target)
	if pageNode == nil {
		logger.Error(errors.New("local webmention target not found: " + webmention.Target))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	page := pageNode.LoadPage()
	page.Webmentions.AddWebmention(webmention)
	page.WriteToFile(true)

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
