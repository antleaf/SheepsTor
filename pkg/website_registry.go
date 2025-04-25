package pkg

type SheepstorWebsiteRegistry struct {
	SourceRoot          string
	WebRoot             string
	GitHubWebHookSecret string
	WebSites            []*SheepstorWebsite
}

func NewSheepstorWebsiteRegistry(sourceRoot, webRoot, gitHubWebHookSecret string) SheepstorWebsiteRegistry {
	r := SheepstorWebsiteRegistry{}
	r.SourceRoot = sourceRoot
	r.WebRoot = webRoot
	r.GitHubWebHookSecret = gitHubWebHookSecret
	r.WebSites = make([]*SheepstorWebsite, 0)
	return r
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
