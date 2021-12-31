package sheepstor

type WebsiteRegistry struct {
	SourceRoot          string
	WebRoot             string
	GitHubWebHookSecret string
	WebSites            []*WebsiteInterface
}

func NewRegistry(sourceRoot, webRoot, gitHubWebHookSecret string) WebsiteRegistry {
	registry := WebsiteRegistry{}
	registry.SourceRoot = sourceRoot
	registry.WebRoot = webRoot
	registry.GitHubWebHookSecret = gitHubWebHookSecret
	registry.WebSites = make([]*WebsiteInterface, 0)
	return registry
}

func (r *WebsiteRegistry) Add(w *WebsiteInterface) {
	r.WebSites = append(r.WebSites, w)
}

func (r *WebsiteRegistry) GetWebsiteByRepoNameAndBranchRef(repoName, branchRef string) *WebsiteInterface {
	for _, wptr := range r.WebSites {
		w := *wptr
		if w.HasRepoNameAndBranchRef(repoName, branchRef) {
			return wptr
		}
	}
	return nil
}

func (r *WebsiteRegistry) GetWebsiteByID(id string) *WebsiteInterface {
	for _, wptr := range r.WebSites {
		w := *wptr
		if w.HasID(id) {
			return wptr
		}
	}
	return nil
}
