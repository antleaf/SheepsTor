package sheepstor

type WebsiteRegistry struct {
	WebSites []*Website
}

var Registry WebsiteRegistry

func InitialiseRegistry(sourceRoot, webRoot string) {
	Registry = WebsiteRegistry{}
	Registry.WebSites = make([]*Website, 0)
	//for _, wc := range websiteConfigs {
	//	if wc.Enabled {
	//		logger.Debug(fmt.Sprintf("Configuring website '%s'", wc.ID))
	//		w := NewWebsite(wc, sourceRootPath, webRoot)
	//		r.WebSites = append(r.WebSites, &w)
	//		logger.Info(fmt.Sprintf("Website '%s' configured OK", wc.ID))
	//	} else {
	//		logger.Warn(fmt.Sprintf("Configuration for website with ID '%s' is not enabled, so not loading this into registry", wc.ID))
	//	}
	//}
}

func (r *WebsiteRegistry) Add(w Website) {
	r.WebSites = append(r.WebSites, &w)
}

func (r *WebsiteRegistry) GetWebsiteByRepoNameAndBranchRef(repoName, branchRef string) *Website {
	for _, w := range r.WebSites {
		if (w.GitRepo.RepoName == repoName) && (w.GitRepo.BranchRef == branchRef) {
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
