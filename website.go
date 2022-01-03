package sheepstor

import (
	"os"
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

func (w *Website) ProvisionSources() error {
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

func (w *Website) CommitAndPush(message string) error {
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
