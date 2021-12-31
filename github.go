package sheepstor

import (
	"github.com/google/go-github/v40/github"
	"net/http"
)

var GitHubWebHookSecret string

func GitHubWebHookHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		payload, err := github.ValidatePayload(req, []byte(GitHubWebHookSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
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
				pushCommitID := *e.HeadCommit.ID
				if localCommitID != pushCommitID {
					err = website.Build()
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
				}
			}
		default:
			return
		}
	}
}
