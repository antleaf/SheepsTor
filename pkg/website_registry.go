package pkg

type SheepstorWebsiteRegistry struct {
	SourceRoot          string
	WebRoot             string
	GitHubWebHookSecret string
	WebSites            []*SheepstorWebsite
}

func NewSheepstorWebsiteRegistry(sourceRoot, webRoot, gitHubWebHookSecret string) SheepstorWebsiteRegistry {
	registry := SheepstorWebsiteRegistry{}
	registry.SourceRoot = sourceRoot
	registry.WebRoot = webRoot
	registry.GitHubWebHookSecret = gitHubWebHookSecret
	registry.WebSites = make([]*SheepstorWebsite, 0)
	return registry
}

func (r *SheepstorWebsiteRegistry) Add(w *SheepstorWebsite) {
	r.WebSites = append(r.WebSites, w)
}

func (r *SheepstorWebsiteRegistry) GetWebsiteByRepoNameAndBranchRef(repoName, branchRef string) *SheepstorWebsite {
	for _, w := range r.WebSites {
		if w.GitRepo.RepoName == repoName && w.GitRepo.BranchRef == branchRef {
			return w
		}
	}
	return nil
}

func (r *SheepstorWebsiteRegistry) GetWebsiteByID(id string) *SheepstorWebsite {
	for _, w := range r.WebSites {
		if w.ID == id {
			return w
		}
	}
	return nil
}
