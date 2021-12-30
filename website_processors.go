package sheepstor

import (
	"fmt"
	"os"
	"os/exec"
)

func HugoProcessor(sourcesPath, targetFolderPathForBuild string) error {
	err := os.MkdirAll(targetFolderPathForBuild, os.ModePerm)
	if err != nil {
		//main.logger.Error(err.Error())
		return err
	}
	hugoCliString := fmt.Sprintf("hugo --quiet --ignoreCache")
	hugoCliString += fmt.Sprintf(" --source %s --destination %s", sourcesPath, targetFolderPathForBuild)
	//main.logger.Debug(fmt.Sprintf("Building website with command '%s'...", hugoCliString))
	hugoCmd := exec.Command("sh", "-c", hugoCliString)
	//var hugoReport []byte
	_, err = hugoCmd.Output()
	//main.logger.Debug(hugoCliString)
	if err != nil {
		//main.logger.Error(string(hugoReport))
		//main.logger.Error(err.Error())
		return err
	}
	return err
}

func DefaultProcessor(sourcesPath, targetFolderPathForBuild string) {
	CopyDir(sourcesPath, targetFolderPathForBuild)
}
