package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type GitRepoConfig struct {
	CloneId       string `yaml:"clone_id"`
	RepoName      string `yaml:"repo_name"`
	BranchName    string `yaml:"branch_name"`
	BranchRef     string `yaml:"-"`
	RepoLocalPath string `yaml:"-"`
}

type SheepsTorProcessorConfig struct {
	BaseURL                   string           `yaml:"base_url"`
	PathProcessors            PathProcessorSet `yaml:"path_processing"`
	WebmentionIoWebhookSecret string           `yaml:"webmention_io_webhook_secret"`
}

type Website struct {
	Id                         string                   `yaml:"id"`
	Enabled                    bool                     `yaml:"enabled"`
	ContentProcessor           string                   `yaml:"content_processor"` //either 'hugo' or nil
	ProcessorRootSubFolderPath string                   `yaml:"processor_root"`    //e.g. a sub-folder in the repo called 'webroot'
	ContentRootSubFolderPath   string                   `yaml:"content_root"`      //for hugo this is 'content' by default
	ProcessorRoot              string                   `yaml:"-"`
	ContentRoot                string                   `yaml:"-"`
	WebRoot                    string                   `yaml:"-"`
	GitRepo                    GitRepoConfig            `yaml:"git"`
	SheepsTorProcessing        SheepsTorProcessorConfig `yaml:"sheepstor"`
	SMap                       Sitemap                  `yaml:"-"`
}

func (w *Website) Configure(sourceRoot, webRoot string) {
	w.GitRepo.BranchRef = fmt.Sprintf("refs/heads/%s", w.GitRepo.BranchName)
	w.GitRepo.RepoLocalPath = filepath.Join(sourceRoot, w.Id)
	if w.ProcessorRootSubFolderPath != "" {
		w.ProcessorRoot = filepath.Join(w.GitRepo.RepoLocalPath, w.ProcessorRootSubFolderPath)
	} else {
		w.ProcessorRoot = w.GitRepo.RepoLocalPath
	}
	if w.ContentRootSubFolderPath != "" {
		w.ContentRoot = filepath.Join(w.ProcessorRoot, w.ContentRootSubFolderPath)
	} else if w.ContentProcessor == "hugo" || w.ContentProcessor == "sheepstor" {
		w.ContentRoot = filepath.Join(w.ProcessorRoot, "content")
	} else {
		w.ContentRoot = w.ProcessorRoot
	}
	if w.ContentProcessor == "sheepstor" {
		for _, pathProcessor := range w.SheepsTorProcessing.PathProcessors.Processors {
			pathProcessor.Initialise(w.SheepsTorProcessing.BaseURL)
		}
		defaultPathProcessor := PathProcessor{Name: "Built-in Default Path Processor", FolderMatchExpression: "(.+)/index\\.md", URLMatchExpression: ""}
		defaultPathProcessor.Initialise(w.SheepsTorProcessing.BaseURL)
		w.SheepsTorProcessing.PathProcessors.DefaultPathProcessor = &defaultPathProcessor
		w.SheepsTorProcessing.PathProcessors.Processors = append(w.SheepsTorProcessing.PathProcessors.Processors, &defaultPathProcessor)
		//w.RegenerateSiteMap()
	}
	w.WebRoot = filepath.Join(webRoot, w.Id)
}

func (w *Website) RegenerateSiteMap() {
	w.SMap = Sitemap{ContentRoot: w.ContentRoot, BaseURL: w.SheepsTorProcessing.BaseURL}
	w.SMap.Build(w.SheepsTorProcessing.PathProcessors)
}

//func (w *Website) GetPermalinkForPage(page Page) string {
//	pagePath := w.IndieWeb.PostsRelativeURLPath
//	pageMonth := strings.Split(page.Metadata.Month, "-")[1]
//	pagePath = strings.Replace(pagePath, ":year", page.Metadata.Year, -1)
//	pagePath = strings.Replace(pagePath, ":month", pageMonth, -1)
//	pagePath = strings.Replace(pagePath, ":slug", page.Slug, -1)
//	return w.IndieWeb.BaseUrl + "/" + pagePath
//}

func (w *Website) Build() error {
	logger.Debug(fmt.Sprintf("Initialising folders for: '%s'....", w.WebRoot))
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
		logger.Error(err.Error())
		return err
	}
	logger.Debug(fmt.Sprintf("Folder: '%s', initialised OK", w.WebRoot))
	switch w.ContentProcessor {
	case "sheepstor":
		err = HugoProcessor(w.ProcessorRoot, targetFolderPathForBuild)
		if err != nil {
			return err
		}
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
			logger.Error(err.Error())
			return err
		}
	} else if os.IsNotExist(err) {
		logger.Debug(fmt.Sprintf("Symlink does not yet exist: '%s'", symbolicLinkPath))
	}
	err = os.Symlink(targetFolderPathForBuild, symbolicLinkPath) //Only switch if successful
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("Built website '%s' OK", w.Id))
	return err
}

func (w *Website) provisionSources() error {
	var err error
	gitFolderPath := filepath.Join(w.GitRepo.RepoLocalPath, ".git")
	if _, err = os.Stat(gitFolderPath); os.IsNotExist(err) {
		logger.Debug(fmt.Sprintf("Git working copy folder does not exist: '%s', creating it now....", w.GitRepo.RepoLocalPath))
		err = os.MkdirAll(w.GitRepo.RepoLocalPath, os.ModePerm)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		logger.Info(fmt.Sprintf("Git working copy folder: '%s', created OK", w.GitRepo.RepoLocalPath))
		err = Clone(w.GitRepo.CloneId, w.GitRepo.BranchRef, w.GitRepo.RepoLocalPath)
		if err != nil {
			logger.Error(err.Error())
			return err
		} else {
			logger.Info(fmt.Sprintf("Cloned sources from '%s' into '%s' OK", w.GitRepo.CloneId, w.GitRepo.RepoLocalPath))
		}
	} else {
		err = Pull(w.GitRepo.RepoLocalPath, w.GitRepo.BranchRef)
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Info(fmt.Sprintf("Sources for website '%s' pulled from '%s' OK", w.Id, w.GitRepo.CloneId))
		}
	}
	return err
}

func (w *Website) CommitAndPush(message string) error {
	err := Pull(w.GitRepo.RepoLocalPath, w.GitRepo.BranchRef)
	if err != nil {
		logger.Error("Git Pull failed " + err.Error())
		return err
	}
	err = CommitAndPush(w.GitRepo.RepoLocalPath, message)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return err
}
