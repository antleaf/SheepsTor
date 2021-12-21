package main

import (
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"io"
	"os"
	"regexp"
	"time"
)

type Sitemap struct {
	ContentRoot      *string
	BaseURL          *string
	PathProcessorSet *PathProcessorSet
	Pages            []*Page
}

func NewSitemap(contentRoot, BaseURL *string, pathProcessorSet *PathProcessorSet) Sitemap {
	var s = Sitemap{
		ContentRoot:      contentRoot,
		BaseURL:          BaseURL,
		PathProcessorSet: pathProcessorSet,
	}
	s.Pages = make([]*Page, 0)
	return s
}

func (s *Sitemap) AddPage(page Page) {
	s.Pages = append(s.Pages, &page)
}

func (s *Sitemap) Build(mediaUploadURLRegex *regexp.Regexp, mediaUploadPath *string) {
	s.Pages = make([]*Page, 0)
	fSys := os.DirFS(*s.ContentRoot)
	logger.Debugf("Building sitemap for content root: %s", *s.ContentRoot)
	paths, err := doublestar.Glob(fSys, "**/index.md")
	if err != nil {
		logger.Error(err.Error())
	}
	for _, path := range paths {
		page := NewPage(path, time.Now(), s.ContentRoot, s.BaseURL, s.PathProcessorSet, mediaUploadURLRegex, mediaUploadPath)
		page.ReadFromFile(false, false)
		s.AddPage(*page)
	}
}

func (s *Sitemap) GetPageByPermalink(permalink string) *Page {
	for _, page := range s.Pages {
		if page.Permalink == permalink {
			return page
		}
	}
	return nil
}

func (s *Sitemap) GetPageByFilePath(filePath string) *Page {
	for _, page := range s.Pages {
		if page.FilePath == filePath {
			return page
		}
	}
	return nil
}

func (s *Sitemap) GetOrCreatePage(filePath string, published time.Time, contentRoot, baseURL *string, pps *PathProcessorSet, mediaUploadURLRegex *regexp.Regexp, mediaUploadPath *string) (*Page, bool) {
	logger.Debugf("Creating page with filepath = %s | contentRoot = %s | baseURL = %s", filePath, *contentRoot, *baseURL)
	newlyCreated := false
	var page *Page
	if filePath != "" {
		page = s.GetPageByFilePath(filePath)
	}
	if page == nil {
		page = NewPage(filePath, published, contentRoot, baseURL, pps, mediaUploadURLRegex, mediaUploadPath)
		newlyCreated = true
	}
	s.AddPage(*page)
	return page, newlyCreated
}

//func (s *Sitemap) GetOrCreatePage(filePath string, pathProcessors PathProcessorSet) (*SitemapNode, bool) {
//	node := s.GetNodeByFilePath(filePath)
//	if node == nil {
//		node = &SitemapNode{FilePath: filePath, ContentRoot: &s.ContentRoot, BaseURL: &s.BaseURL}
//		pathProcessors.AssignPathProcessorToSitemapNode(node)
//		node.ExtrapolatePermalink()
//		s.Nodes = append(s.Nodes, node)
//		return node, true
//	}
//	return node, false
//}

func (s *Sitemap) Dump(w io.Writer) {
	for _, p := range s.Pages {
		line := fmt.Sprintf("%s => %s [%s]\n", p.FilePath, p.Permalink, p.PathProcessor.Name)
		w.Write([]byte(line))
	}
}
