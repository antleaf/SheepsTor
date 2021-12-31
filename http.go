package sheepstor

import (
	"github.com/google/go-github/v40/github"
	"net/http"
)

func (r *WebsiteRegistry) GitHubWebHookHandler(resp http.ResponseWriter, req *http.Request) {
	logger.Debug("Handling GitHUb webhook post....")
	if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(r.GitHubWebHookSecret))
		if err != nil {
			logger.Error(err.Error())
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
			logger.Error(err.Error())
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		switch e := event.(type) {
		case *github.PushEvent:
			logger.Debug("Github push event received")
			websitePtr := r.GetWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
			if websitePtr != nil {
				website := *websitePtr
				logger.Debugf("Website identified from GitHub push event; '%s'", website.GetID())
				gitRepo := website.GetGitRepo()
				localCommitID := gitRepo.GetHeadCommitID()
				pushCommitID := *e.HeadCommit.ID
				if localCommitID != pushCommitID {
					logger.Debugf("Attempting to build website '%s'", website.GetID())
					err = website.ProvisionSources()
					if err != nil {
						logger.Error(err.Error())
						http.Error(resp, err.Error(), http.StatusBadRequest)
						return
					}
					err = website.Build()
					if err != nil {
						logger.Error(err.Error())
						http.Error(resp, err.Error(), http.StatusBadRequest)
						return
					} else {
						logger.Infof("Built website '%s'", website.GetID())
					}
				}
			} else {
				logger.Errorf("Website with repo name '%s' and branch ref '%s' not found", e.GetRepo().GetFullName(), e.GetRef())
			}
		default:
			return
		}
	}
}
