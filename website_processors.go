package sheepstor

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
