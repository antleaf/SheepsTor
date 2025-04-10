package internal

import (
	"github.com/google/go-github/v40/github"
	"net/http"
)

type WebsiteRegistry struct {
	SourceRoot          string
	WebRoot             string
	GitHubWebHookSecret string
	WebSites            []*Website
}

func NewRegistry(sourceRoot, webRoot, gitHubWebHookSecret string) WebsiteRegistry {
	registry := WebsiteRegistry{}
	registry.SourceRoot = sourceRoot
	registry.WebRoot = webRoot
	registry.GitHubWebHookSecret = gitHubWebHookSecret
	registry.WebSites = make([]*Website, 0)
	return registry
}

func (r *WebsiteRegistry) Add(w *Website) {
	r.WebSites = append(r.WebSites, w)
}

func (r *WebsiteRegistry) GetWebsiteByRepoNameAndBranchRef(repoName, branchRef string) *Website {
	for _, w := range r.WebSites {
		if w.GitRepo.RepoName == repoName && w.GitRepo.BranchRef == branchRef {
			return w
		}
	}
	return nil
}

func (r *WebsiteRegistry) GetWebsiteByID(id string) *Website {
	for _, w := range r.WebSites {
		if w.ID == id {
			return w
		}
	}
	return nil
}

func (r *WebsiteRegistry) GitHubWebHookHandler(resp http.ResponseWriter, req *http.Request) {
	Log.Debug("Handling GitHUb webhook post....")
	if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(r.GitHubWebHookSecret))
		if err != nil {
			Log.Error(err.Error())
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
			Log.Error(err.Error())
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		switch e := event.(type) {
		case *github.PushEvent:
			Log.Debug("Github push event received")
			websitePtr := r.GetWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
			if websitePtr != nil {
				website := *websitePtr
				Log.Debugf("Website identified from GitHub push event; '%s'", website.ID)
				gitRepo := website.GitRepo
				localCommitID := gitRepo.GetHeadCommitID()
				pushCommitID := *e.HeadCommit.ID
				if localCommitID != pushCommitID {
					Log.Debugf("Attempting to build website '%s'", website.ID)
					err = website.ProvisionSources()
					if err != nil {
						Log.Error(err.Error())
						http.Error(resp, err.Error(), http.StatusBadRequest)
						return
					}
					err = website.Build()
					if err != nil {
						Log.Error(err.Error())
						http.Error(resp, err.Error(), http.StatusBadRequest)
						return
					} else {
						Log.Infof("Built website '%s'", website.ID)
					}
				}
			} else {
				Log.Errorf("Website with repo name '%s' and branch ref '%s' not found", e.GetRepo().GetFullName(), e.GetRef())
			}
		default:
			return
		}
	}
}
