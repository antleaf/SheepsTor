package internal

import (
	"fmt"
	"os"
	"os/exec"
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

func IndexForSearch(targetFolderPathForBuild string) error {
	var err error
	indexCmdString := fmt.Sprintf("pagefind --site %s", targetFolderPathForBuild)
	indexCmd := exec.Command("sh", "-c", indexCmdString)
	_, err = indexCmd.Output()
	if err != nil {
		return err
	}
	return err
}
