package main

import "regexp"

type PathProcessor struct {
	Name                  string         `yaml:"name"`
	FolderMatchExpression string         `yaml:"folder_match_expression"`
	UrlGenerationPattern  string         `yaml:"url_generation_pattern"`
	URLMatchExpression    string         `yaml:"url_match_expression"`
	FolderRegex           *regexp.Regexp `yaml:"-"`
	URLRegex              *regexp.Regexp `yaml:"-"`
	BaseURL               string         `yaml:"-"` //comes from website config
}

type PathProcessorSet struct {
	Processors           []*PathProcessor `yaml:"processors"`
	DefaultPathProcessor *PathProcessor   `yaml:"-"`
}

func (pp *PathProcessor) Initialise(baseURL string) {
	pp.BaseURL = baseURL
	pp.FolderRegex = regexp.MustCompile(pp.FolderMatchExpression)
	pp.URLRegex = regexp.MustCompile(pp.URLMatchExpression)
}

func (pps *PathProcessorSet) AssignPathProcessorToSitemapNode(node *SitemapNode) {
	if len(node.FilePath) > 0 {
		for _, pp := range pps.Processors {
			if len(pp.FolderRegex.FindAllStringSubmatch(node.FilePath, -1)) > 0 {
				//logger.Infof("Assigned path %s to %s", node.FilePath, pp.Name)
				node.PathProcessor = pp
				return
			}
		}
		logger.Warnf("Did not find path processor for file path %s", node.FilePath)
		node.PathProcessor = pps.DefaultPathProcessor
	} else {
		logger.Errorf("No file path set for node, unable to match path processor, setting default path processor")
		node.PathProcessor = pps.DefaultPathProcessor
	}
}
