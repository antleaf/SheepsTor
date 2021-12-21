package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func HugoProcessor(sourcesPath, targetFolderPathForBuild string) error {
	err := os.MkdirAll(targetFolderPathForBuild, os.ModePerm)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	hugoCliString := fmt.Sprintf("hugo --quiet --ignoreCache")
	hugoCliString += fmt.Sprintf(" --source %s --destination %s", sourcesPath, targetFolderPathForBuild)
	logger.Debug(fmt.Sprintf("Building website with command '%s'...", hugoCliString))
	hugoCmd := exec.Command("sh", "-c", hugoCliString)
	var hugoReport []byte
	hugoReport, err = hugoCmd.Output()
	logger.Debug(hugoCliString)
	if err != nil {
		logger.Error(string(hugoReport))
		logger.Error(err.Error())
		return err
	}
	return err
}

func DefaultProcessor(sourcesPath, targetFolderPathForBuild string) {
	CopyDir(sourcesPath, targetFolderPathForBuild)
}

func ProcessWebsite(website *Website) {
	err := website.provisionSources()
	if err != nil {
		logger.Error(err.Error())
	} else {
		err = website.Build()
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func processWebsiteInSynchronousWorker(website *Website, wg *sync.WaitGroup) {
	err := website.provisionSources()
	if err == nil {
		err = website.Build()
		if err != nil {
			logger.Error(err.Error())
		}
	}
	wg.Done()
}

func ProcessAllWebsites() {
	var wg sync.WaitGroup
	for _, website := range registry.WebSites {
		wg.Add(1)
		go processWebsiteInSynchronousWorker(website, &wg)
	}
	wg.Wait()
	systemReadyToAcceptUpdateRequests = true
}
