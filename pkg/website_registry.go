package pkg

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
