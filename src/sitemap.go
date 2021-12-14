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
	page.FilePath = filepath.Join(*sn.ContentRoot, sn.FilePath)
	page.Permalink = sn.Permalink
	page.BaseURL = sn.BaseURL
	page.ReadFromFile()
	return page
}

func (sn *SitemapNode) ExtrapolatePermalink() {
	permalink := sn.PathProcessor.FolderRegex.ReplaceAllString(sn.FilePath, sn.PathProcessor.UrlGenerationPattern)
	page := sn.LoadPage()
	permalink = strings.Replace(permalink, "{year}", page.Metadata.Year(), -1)
	permalink = strings.Replace(permalink, "{month}", page.Metadata.Month(), -1)
	permalink = strings.Replace(permalink, "{slug}", page.Slug, -1)
	sn.Permalink = *sn.BaseURL + "/" + permalink
}

func (s *Sitemap) GetNodeByPermalink(permalink string) *SitemapNode {
	for _, node := range s.Nodes {
		if node.Permalink == permalink {
			return node
		}
	}
	return nil
}

//func (s *Sitemap) GetNodeByPermalink(permalink string) (*SitemapNode, error) {
//	var err error
//	for _, node := range s.Nodes {
//		if node.Permalink == permalink {
//			return node, err
//		}
//	}
//	err = errors.New("sitemap node with permalink " + permalink + " not found")
//	return &SitemapNode{}, err
//}
