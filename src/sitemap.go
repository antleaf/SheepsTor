package main

import (
	"github.com/bmatcuk/doublestar/v4"
	"os"
	"path/filepath"
	"strings"
)

type SitemapNode struct {
	FilePath      string
	Permalink     string
	PathProcessor *PathProcessor
	ContentRoot   *string
	BaseURL       *string
}

type Sitemap struct {
	Nodes       []*SitemapNode
	ContentRoot string
	BaseURL     string
}

func (s *Sitemap) Build(pathProcessors PathProcessorSet) {
	s.Nodes = make([]*SitemapNode, 0)
	fsys := os.DirFS(s.ContentRoot)
	logger.Debug(s.ContentRoot)
	paths, err := doublestar.Glob(fsys, "**/index.md")
	if err != nil {
		logger.Error(err.Error())
	}
	for _, path := range paths {
		sitemapNode := SitemapNode{FilePath: path, ContentRoot: &s.ContentRoot, BaseURL: &s.BaseURL}
		pathProcessors.AssignPathProcessorToSitemapNode(&sitemapNode)
		sitemapNode.ExtrapolatePermalink()
		s.Nodes = append(s.Nodes, &sitemapNode)
	}
}

func (sn *SitemapNode) LoadPage() Page {
	page := Page{}
	page.ReadFromFile(filepath.Join(*sn.ContentRoot, sn.FilePath))
	return page
}

func (sn *SitemapNode) ExtrapolatePermalink() {
	permalink := sn.PathProcessor.FolderRegex.ReplaceAllString(sn.FilePath, sn.PathProcessor.UrlGenerationPattern)
	page := sn.LoadPage()
	pageMonth := ""
	pageYearMonthStringSlice := strings.Split(page.Metadata.Month, "-")
	if len(pageYearMonthStringSlice) == 2 {
		pageMonth = strings.Split(page.Metadata.Month, "-")[1]
	}
	permalink = strings.Replace(permalink, "{year}", page.Metadata.Year, -1)
	permalink = strings.Replace(permalink, "{month}", pageMonth, -1)
	permalink = strings.Replace(permalink, "{slug}", page.Slug, -1)
	sn.Permalink = *sn.BaseURL + "/" + permalink
}