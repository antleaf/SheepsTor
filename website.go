package sheepstor

import (
	"path/filepath"
)

type Website struct {
	ID               string
	ContentProcessor string //either 'hugo' or nil
	ProcessorRoot    string
	WebRoot          string
	GitRepo          GitRepo
}

func NewWebsite(id, contentProcessor, processorRoot, sourceRoot, webRoot, repoCloneID, repoName, repoBranchName string) Website {
	var w = Website{
		ID:               id,
		ContentProcessor: contentProcessor,
	}
	w.WebRoot = filepath.Join(webRoot, w.ID)
	w.GitRepo = NewGitRepo(repoCloneID, repoName, repoBranchName, filepath.Join(sourceRoot, w.ID))
	if processorRoot != "" {
		w.ProcessorRoot = filepath.Join(w.GitRepo.RepoLocalPath, processorRoot)
	} else {
		w.ProcessorRoot = w.GitRepo.RepoLocalPath
	}
	return w
}

func (w *Website) Build() error {
	return Build(w)
}

func (w *Website) ProvisionSources() error {
	return ProvisionSources(w)
}

func (w *Website) CommitAndPush(message string) error {
	return CommitAndPush(w, message)
}

func (w *Website) HasID(id string) bool {
	return w.ID == id
}

func (w *Website) HasRepoNameAndBranchRef(repoName, branchRef string) bool {
	return w.GitRepo.RepoName == repoName && w.GitRepo.BranchRef == branchRef
}

func (w *Website) GetID() string {
	return w.ID
}

func (w *Website) GetGitRepo() GitRepo {
	return w.GitRepo
}

func (w *Website) GetWebRoot() string {
	return w.WebRoot
}

func (w *Website) GetContentProcessor() string {
	return w.ContentProcessor
}

func (w *Website) GetProcessorRoot() string {
	return w.ProcessorRoot
}
