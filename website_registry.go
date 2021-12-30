package sheepstor

type WebsiteRegistry struct {
	WebSites []*WebsiteInterface
}

var Registry WebsiteRegistry

func InitialiseRegistry(sourceRoot, webRoot string) {
	Registry = WebsiteRegistry{}
	Registry.WebSites = make([]*WebsiteInterface, 0)
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
