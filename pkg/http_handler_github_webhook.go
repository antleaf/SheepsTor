package pkg

import (
	"github.com/google/go-github/v40/github"
	"net/http"
)

func GitHubWebHookHandler(resp http.ResponseWriter, req *http.Request) {
	log.Debug("Handling GitHUb webhook post....")
	if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(Registry.GitHubWebHookSecret))
		if err != nil {
			log.Error(err.Error())
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
			log.Error(err.Error())
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		switch e := event.(type) {
		case *github.PushEvent:
			log.Debug("Github push event received")
			websitePtr := Registry.GetWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
			if websitePtr != nil {
				website := *websitePtr
				log.Debugf("Website identified from GitHub push event; '%s'", website.ID)
				gitRepo := website.GitRepo
				localCommitID := gitRepo.GetHeadCommitID()
				pushCommitID := *e.HeadCommit.ID
				if localCommitID != pushCommitID {
					log.Debugf("Attempting to build website '%s'", website.ID)
					err = website.ProvisionSources()
					if err != nil {
						log.Error(err.Error())
						http.Error(resp, err.Error(), http.StatusBadRequest)
						return
					}
					err = website.Build()
					if err != nil {
						log.Error(err.Error())
						http.Error(resp, err.Error(), http.StatusBadRequest)
						return
					} else {
						log.Infof("Built website '%s'", website.ID)
					}
				}
			} else {
				log.Errorf("Website with repo name '%s' and branch ref '%s' not found", e.GetRepo().GetFullName(), e.GetRef())
			}
		default:
			return
		}
	}
}
