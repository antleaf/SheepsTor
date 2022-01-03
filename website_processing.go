package sheepstor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func HugoProcessor(sourcesPath, targetFolderPathForBuild string) error {
	err := os.MkdirAll(targetFolderPathForBuild, os.ModePerm)
	if err != nil {
		return err
	}
	hugoCliString := fmt.Sprintf("hugo --quiet --ignoreCache")
	hugoCliString += fmt.Sprintf(" --source %s --destination %s", sourcesPath, targetFolderPathForBuild)
	hugoCmd := exec.Command("sh", "-c", hugoCliString)
	_, err = hugoCmd.Output()
	if err != nil {
		return err
	}
	return err
}

func DefaultProcessor(sourcesPath, targetFolderPathForBuild string) {
	CopyDir(sourcesPath, targetFolderPathForBuild)
}

func ProvisionSources(w WebsiteInterface) error {
	var err error
	gitRepo := w.GetGitRepo()
	gitFolderPath := filepath.Join(gitRepo.RepoLocalPath, ".git")
	if _, err = os.Stat(gitFolderPath); os.IsNotExist(err) {
		err = os.MkdirAll(gitRepo.RepoLocalPath, os.ModePerm)
		if err != nil {
			return err
		}
		err = gitRepo.Clone()
		if err != nil {
			return err
		}
	} else {
		err = gitRepo.Pull()
		if err != nil {
			return err
		}
	}
	return err
}

func CommitAndPush(w WebsiteInterface, message string) error {
	gitRepo := w.GetGitRepo()
	err := gitRepo.Pull()
	if err != nil {
		return err
	}
	err = gitRepo.CommitAndPush(message)
	if err != nil {
		return err
	}
	return err
}

func Build(w WebsiteInterface) error {
	var err error
	targetFolderPathForBuild := filepath.Join(w.GetWebRoot(), "public_1")
	symbolicLinkPath := filepath.Join(w.GetWebRoot(), "public")
	currentSymLinkTargetPath, readlinkErr := os.Readlink(symbolicLinkPath)
	if readlinkErr == nil {
		if currentSymLinkTargetPath == filepath.Join(w.GetWebRoot(), "public_1") {
			targetFolderPathForBuild = filepath.Join(w.GetWebRoot(), "public_2")
		}
	}
	if _, statErr := os.Stat(targetFolderPathForBuild); statErr == nil {
		os.RemoveAll(targetFolderPathForBuild)
	}
	err = os.MkdirAll(w.GetWebRoot(), os.ModePerm)
	err = os.MkdirAll(filepath.Join(w.GetWebRoot(), "logs"), os.ModePerm)
	if err != nil {
		return err
	}
	switch w.GetContentProcessor() {
	case "hugo":
		err = HugoProcessor(w.GetProcessorRoot(), targetFolderPathForBuild)
		if err != nil {
			return err
		}
	default:
		DefaultProcessor(w.GetProcessorRoot(), targetFolderPathForBuild)
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
