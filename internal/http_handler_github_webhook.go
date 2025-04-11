package internal

import (
	"github.com/google/go-github/v40/github"
	"net/http"
)

func GitHubWebHookHandler(resp http.ResponseWriter, req *http.Request) {
	Log.Debug("Handling GitHUb webhook post....")
	if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(Registry.GitHubWebHookSecret))
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
			websitePtr := Registry.GetWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
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
