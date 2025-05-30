package pkg

import (
	toolbox2go "github.com/antleaf/toolbox2go"
	"os"
	"path/filepath"
)

type SheepstorWebsite struct {
	ID               string
	ContentProcessor string //either 'hugo' or nil
	ProcessorRoot    string
	WebRoot          string
	IndexForSearch   bool
	GitRepo          toolbox2go.GitRepo
}

func NewSheepstorWebsite(id, contentProcessor, processorRoot, sourceRoot, webRoot, repoCloneID, repoName, repoBranchName string, indexForSearch bool) SheepstorWebsite {
	var w = SheepstorWebsite{
		ID:               id,
		ContentProcessor: contentProcessor,
		IndexForSearch:   indexForSearch,
	}
	w.WebRoot = filepath.Join(webRoot, w.ID)
	w.GitRepo = toolbox2go.NewGitRepo(repoCloneID, repoName, repoBranchName, filepath.Join(sourceRoot, w.ID))
	if processorRoot != "" {
		w.ProcessorRoot = filepath.Join(w.GitRepo.RepoLocalPath, processorRoot)
	} else {
		w.ProcessorRoot = w.GitRepo.RepoLocalPath
	}
	return w
}

func (w *SheepstorWebsite) Build() error {
	var err error
	targetFolderPathForBuild := filepath.Join(w.WebRoot, "public_1")
	symbolicLinkPath := filepath.Join(w.WebRoot, "public")
	currentSymLinkTargetPath, readlinkErr := os.Readlink(symbolicLinkPath)
	if readlinkErr == nil {
		if currentSymLinkTargetPath == filepath.Join(w.WebRoot, "public_1") {
			targetFolderPathForBuild = filepath.Join(w.WebRoot, "public_2")
		}
	}
	if _, statErr := os.Stat(targetFolderPathForBuild); statErr == nil {
		os.RemoveAll(targetFolderPathForBuild)
	}
	err = os.MkdirAll(w.WebRoot, os.ModePerm)
	err = os.MkdirAll(filepath.Join(w.WebRoot, "logs"), os.ModePerm)
	if err != nil {
		return err
	}
	switch w.ContentProcessor {
	case "hugo":
		err = HugoProcessor(w.ProcessorRoot, targetFolderPathForBuild)
		if err != nil {
			return err
		}
	default:
		DefaultProcessor(w.ProcessorRoot, targetFolderPathForBuild)
	}
	if w.IndexForSearch {
		err = IndexForSearch(targetFolderPathForBuild)
		if err != nil {
			return err
		}
	}
	if _, err = os.Lstat(symbolicLinkPath); err == nil {
		if err = os.Remove(symbolicLinkPath); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		//do nothing?
	}
	err = os.Symlink(targetFolderPathForBuild, symbolicLinkPath) //Only switch if successful
	if err != nil {
		return err
	}
	return err
}

func (w *SheepstorWebsite) ProvisionSources() error {
	var err error
	gitFolderPath := filepath.Join(w.GitRepo.RepoLocalPath, ".git")
	if _, err = os.Stat(gitFolderPath); os.IsNotExist(err) {
		err = os.MkdirAll(w.GitRepo.RepoLocalPath, os.ModePerm)
		if err != nil {
			return err
		}
		err = w.GitRepo.Clone()
		if err != nil {
			return err
		}
	} else {
		err = w.GitRepo.Pull()
		if err != nil {
			return err
		}
	}
	return err
}

func (w *SheepstorWebsite) CommitAndPush(message string) error {
	err := w.GitRepo.Pull()
	if err != nil {
		return err
	}
	err = w.GitRepo.CommitAndPush(message)
	if err != nil {
		return err
	}
	return err
}
