package main

import "fmt"

type WebsiteRegistry struct {
	WebSites []*Website
}

func NewRegistry(websiteConfigs []WebsiteConfig, sourceRootPath, webRoot string) WebsiteRegistry {
	var r WebsiteRegistry
	for _, wc := range websiteConfigs {
		if wc.Enabled {
			logger.Debug(fmt.Sprintf("Configuring website '%s'", wc.ID))
			w := NewWebsite(wc, sourceRootPath, webRoot)
			r.WebSites = append(r.WebSites, &w)
			logger.Info(fmt.Sprintf("Website '%s' configured OK", wc.ID))
		} else {
			logger.Warn(fmt.Sprintf("Configuration for website with ID '%s' is not enabled, so not loading this into registry", wc.ID))
		}
	}
	return r
}

func (r *WebsiteRegistry) getWebsiteByRepoNameAndBranchRef(repoName, branchRef string) *Website {
	for _, w := range r.WebSites {
		if (w.GitRepo.RepoName == repoName) && (w.GitRepo.BranchRef == branchRef) {
			return w
		}
	}
	return nil
}

func (r *WebsiteRegistry) getWebsiteByID(id string) *Website {
	for _, w := range r.WebSites {
		if w.ID == id {
			return w
		}
	}
	return nil
}
