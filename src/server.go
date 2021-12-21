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
	//r.Get("/micropub-discovery/{websiteID}", MicroPubDiscoveryHandler)
	r.Route("/micropub/{websiteID}", func(r chi.Router) {
		r.Use(MicroPubAuthorisationMiddleware)
		r.Get("/", MicroPubGetHandler)
		r.Post("/", MicroPubPostHandler)
		r.Post("/media", MicroPubMediaHandler)
	})
	return r
}

func WebMentionIOHookHandler(w http.ResponseWriter, r *http.Request) {
	websiteID := chi.URLParam(r, "websiteID")
	logger.Debugf("Received webMention.io webhook for website with ID: %s", websiteID)
	website := registry.getWebsiteByID(websiteID)
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
	webMentionPayload := WebMentionIOPayload{}
	webMention, err := webMentionPayload.LoadAndValidate(payloadJson, website.IndieWeb.WebMentionIoWebhookSecret)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Debugf("Processing incoming webMention with source: %s & target: %s...", webMention.Source, webMention.Target)
	pageFullFilePath := website.GetPagePathForPermalink(webMention.Target)
	page, err := website.LoadPage(pageFullFilePath) //TODO: GetPagePathForPermalink() needs to be developed
	if err != nil {
		logger.Error(errors.New("local webMention target not found: " + webMention.Target))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	page.WebMentions.AddWebMention(webMention)
	page.WriteToFile(pageFullFilePath)
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
		payload, err := github.ValidatePayload(req, []byte(os.Getenv(config.GitHubWebHookSecretEnvKey)))
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
			websitePtr := registry.getWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
			if websitePtr != nil {
				localCommitID := websitePtr.GitRepo.GetHeadCommitID()
				logger.Debugf("Local commit ID = %s", localCommitID)
				pushCommitID := *e.HeadCommit.ID
				logger.Debugf("Head commit ID from push event = %v", pushCommitID)
				if localCommitID != pushCommitID {
					ProcessWebsite(websitePtr)
				} else {
					logger.Debugf("Local and Push commit IDs are the same - no need to rebuild the site")
				}
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
