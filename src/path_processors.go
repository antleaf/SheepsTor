package main

import "regexp"

type PathProcessor struct {
	Name                  string
	FolderMatchExpression string
	UrlGenerationPattern  string
	FileGenerationPattern string
	FolderRegex           *regexp.Regexp
	BaseURL               string
}

type PathProcessorSet struct {
	DefaultPathProcessor PathProcessor
	PathProcessors       []PathProcessor
}

func NewPathProcessorSet(defaultPPConfig PathProcessorConfig, ppConfigs []PathProcessorConfig, baseURL string) PathProcessorSet {
	var pps PathProcessorSet
	pps.DefaultPathProcessor = NewPathProcessor(defaultPPConfig, baseURL)
	for _, ppConfig := range ppConfigs {
		pps.PathProcessors = append(pps.PathProcessors, NewPathProcessor(ppConfig, baseURL))
	}
	return pps
}

var DefaultPPConfig = PathProcessorConfig{
	Name:                  "Built-in Default Path Processor",
	FolderMatchExpression: "([a-zA-Z0-9_\\/-]*)/index\\.md$",
	UrlGenerationPattern:  "$1/",
	FileGenerationPattern: "/{slug}/index.md",
}

func NewPathProcessor(ppConfig PathProcessorConfig, baseURL string) PathProcessor {
	var pp = PathProcessor{
		Name:                  ppConfig.Name,
		FolderMatchExpression: ppConfig.FolderMatchExpression,
		UrlGenerationPattern:  ppConfig.UrlGenerationPattern,
		FileGenerationPattern: ppConfig.FileGenerationPattern,
		FolderRegex:           nil,
		BaseURL:               baseURL,
	}
	pp.FolderRegex = regexp.MustCompile(pp.FolderMatchExpression)
	return pp
}

func (pps *PathProcessorSet) SelectPathProcessorForPath(path string) *PathProcessor {
	if len(path) > 0 {
		for _, pp := range pps.PathProcessors {
			if len(pp.FolderRegex.FindAllStringSubmatch(path, -1)) > 0 {
				return &pp
			}
		}
		logger.Warnf("Did not find path processor for file path %s so using default processor", path)
	} else {
		logger.Errorf("No file path set for node, unable to match path processor")
	}
	return &pps.DefaultPathProcessor
}

//func (pps *PathProcessorSet) GetPathProcessorByName(name string) *PathProcessor {
//	for _, pp := range pps.PathProcessors {
//		if pp.Name == name {
//			return &pp
//		}
//	}
//	return &pps.DefaultPathProcessor
//}
