package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Website struct {
	ID                   string
	ContentProcessor     string //either 'hugo' or nil
	ProcessorRoot        string
	ContentRoot          string
	WebRoot              string
	BaseURL              string
	GitRepo              GitRepo
	IndieWeb             IndieWeb
	PathProcessorSet     PathProcessorSet
	PathProcessors       []PathProcessor
	DefaultPathProcessor PathProcessor
	SiteMap              Sitemap
}

func NewWebsite(wConfig WebsiteConfig, sourceRootPath, webRoot string) Website {
	var w = Website{
		ID:               wConfig.ID,
		ContentProcessor: wConfig.ContentProcessor,
	}
	w.WebRoot = filepath.Join(webRoot, w.ID)
	w.GitRepo = NewGitRepo(wConfig.GitRepoConfig, filepath.Join(sourceRootPath, w.ID))
	if wConfig.ProcessorRootSubFolderPath != "" {
		w.ProcessorRoot = filepath.Join(w.GitRepo.RepoLocalPath, wConfig.ProcessorRootSubFolderPath)
	} else {
		w.ProcessorRoot = w.GitRepo.RepoLocalPath
	}
	if wConfig.ContentRootSubFolderPath != "" {
		w.ContentRoot = filepath.Join(w.ProcessorRoot, wConfig.ContentRootSubFolderPath)
	} else if w.ContentProcessor == "hugo" || w.ContentProcessor == "sheepstor" {
		w.ContentRoot = filepath.Join(w.ProcessorRoot, "content")
	} else {
		w.ContentRoot = w.ProcessorRoot
	}
	if w.ContentProcessor == "sheepstor" {
		w.BaseURL = wConfig.SheepsTorProcessing.BaseURL
		w.PathProcessorSet = NewPathProcessorSet(DefaultPPConfig, wConfig.SheepsTorProcessing.PathProcessorConfigs, w.BaseURL)
		w.IndieWeb = NewIndieWeb(wConfig.SheepsTorProcessing.IndieWebConfig, w.BaseURL, w.ProcessorRoot)
		w.SiteMap = NewSitemap(&w.ContentRoot, &w.BaseURL, &w.PathProcessorSet)
		w.SiteMap.Build(w.IndieWeb.MediaUploadURLRegex, &w.IndieWeb.MediaUploadPath)
	}
	return w
}

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
	logger.Info(fmt.Sprintf("Built website '%s' OK", w.ID))
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
			logger.Info(fmt.Sprintf("Sources for website '%s' pulled from '%s' OK", w.ID, w.GitRepo.CloneId))
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

func (w *Website) StoreTempMediaFileAndReturnURL(mediaFile multipart.File, fileName string) (string, string, error) {
	mediaFilePath := filepath.Join(w.IndieWeb.MediaUploadPath, fileName)
	mediaFileURL := fmt.Sprintf("%s/%s", w.IndieWeb.MediaUploadBaseURL, fileName)
	defer mediaFile.Close()
	raw, err := ioutil.ReadAll(mediaFile)
	if err != nil {
		logger.Error(err.Error())
		return mediaFilePath, mediaFileURL, err
	}
	err = os.WriteFile(mediaFilePath, raw, os.ModePerm)
	return mediaFilePath, mediaFileURL, err
}
