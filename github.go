package sheepstor

import (
	"github.com/google/go-github/v40/github"
	"net/http"
)

var GitHubWebHookSecret string

func GitHubWebHookHandler(w http.ResponseWriter, req *http.Request) {
	//main.logger.Debug("Handling GitHUb webhook post....")
	if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(GitHubWebHookSecret))
		if err != nil {
			//main.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
			//main.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch e := event.(type) {
		case *github.PushEvent:
			websitePtr := Registry.GetWebsiteByRepoNameAndBranchRef(e.GetRepo().GetFullName(), e.GetRef())
			if websitePtr != nil {
				website := *websitePtr
				gitRepo := website.GetGitRepo()
				localCommitID := gitRepo.GetHeadCommitID()
				//main.logger.Debugf("Local commit ID = %s", localCommitID)
				pushCommitID := *e.HeadCommit.ID
				//main.logger.Debugf("Head commit ID from push event = %v", pushCommitID)
				if localCommitID != pushCommitID {
					err = website.Build()
					if err != nil {
						//main.logger.Error(err.Error())
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
				} else {
					//main.logger.Debugf("Local and Push commit IDs are the same - no need to rebuild the site")
				}
			}
		default:
			//main.logger.Error(fmt.Sprintf("Not a push event %s", github.WebHookType(req)))
			return
		}
	}
}
