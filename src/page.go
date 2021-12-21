package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Page struct {
	Permalink           string
	FilePath            string
	Title               string
	ContentRoot         *string
	Published           time.Time
	Slug                string
	Metadata            FrontMatter
	Content             string
	WebMentions         *WebMentionSet
	BaseURL             *string
	PathProcessor       *PathProcessor
	MediaUploadURLRegex *regexp.Regexp
	MediaUploadPath     *string
}

func NewPage(filePath string, published time.Time, contentRoot, baseURL *string, pps *PathProcessorSet, mediaUploadURLRegex *regexp.Regexp, mediaUploadPath *string) *Page {
	page := Page{
		ContentRoot:         contentRoot,
		FilePath:            filePath,
		BaseURL:             baseURL,
		Published:           published,
		MediaUploadURLRegex: mediaUploadURLRegex,
		MediaUploadPath:     mediaUploadPath,
	}
	page.Title = strings.ReplaceAll(page.Slug, "_", " ")
	page.Title = strings.Title(page.Title)
	page.Slug = filepath.Base(filepath.Dir(page.FilePath))
	//if page.FilePath != "" && page.Slug == "" {
	//	page.Slug = filepath.Base(filepath.Dir(page.FilePath))
	//} else if page.FilePath == "" && page.Slug != "" {
	//	//page.FilePath = filepath.Join(*contentRoot, page.Slug, "index.md")
	//}
	page.PathProcessor = pps.SelectPathProcessorForPath(page.FilePath)
	page.ExtrapolatePermalink()
	webMentionSet := NewWebMentionSet(filepath.Join(page.FullFolderPath(), "webmentions.csv"))
	page.WebMentions = &webMentionSet
	return &page
}

func (p *Page) FullFilePath() string {
	return filepath.Join(*p.ContentRoot, p.FilePath)
}

func (p *Page) FullFolderPath() string {
	return filepath.Join(*p.ContentRoot, filepath.Dir(p.FilePath))
}

func (p *Page) WriteToFile(withMetadata, writeWebMentions bool) error {
	var fullPostContent string
	var err error
	if withMetadata {
		fullPostContent, err = p.Metadata.ToYaml()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	fullPostContent = fullPostContent + p.Content
	os.MkdirAll(p.FullFolderPath(), os.ModePerm)
	err = ioutil.WriteFile(p.FullFilePath(), []byte(fullPostContent), os.ModePerm)
	if writeWebMentions {
		p.ExtractLinksForWebMentions()
		p.WebMentions.SaveToFile()
	}

	for _, match := range p.MediaUploadURLRegex.FindAllStringSubmatch(p.Content, -1) {
		mediaFileName := match[2]
		logger.Debugf("Found media with file name %s", mediaFileName)
		oldFilePath := filepath.Join(*p.MediaUploadPath, mediaFileName)
		newFilePath := filepath.Join(p.FullFolderPath(), mediaFileName)
		err = os.Rename(oldFilePath, newFilePath)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	fullPostContent = fullPostContent + p.MediaUploadURLRegex.ReplaceAllString(p.Content, "./$2")
	return err
}

func (p *Page) ReadFromFile(loadContent, loadWebMentions bool) error {
	var err error
	fileBytes, err := ioutil.ReadFile(p.FullFilePath())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = p.ReadFromString(string(fileBytes), loadContent)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if loadWebMentions {
		p.WebMentions.LoadFromFile()
	}
	return err
}

func (p *Page) ReadFromString(rawPage string, loadContent bool) error {
	var err error
	slug := p.FilePath[:len(p.FilePath)-8]
	slug = filepath.Base(slug)
	p.Slug = slug
	metadata, content, err := ParseTextForFrontMatterAndContent(rawPage)
	if err == nil {
		p.Metadata = *metadata
		p.Published = p.Metadata.Date
		p.Title = p.Metadata.Title
	}
	if loadContent {
		p.Content = content
	}
	return err
}

func (p *Page) ExtrapolatePermalink() {
	permalink := p.PathProcessor.FolderRegex.ReplaceAllString(p.FilePath, p.PathProcessor.UrlGenerationPattern)
	permalink = strings.Replace(permalink, "{year}", p.Published.Format("2006"), -1)
	permalink = strings.Replace(permalink, "{month}", p.Published.Format("01"), -1)
	permalink = strings.Replace(permalink, "{slug}", p.Slug, -1)
	p.Permalink = *p.BaseURL + "/" + permalink
}

func (p *Page) ExtractLinksForWebMentions() {
	targets := ExtractLinkURLs(p.Content, *p.BaseURL)
	for _, target := range targets {
		p.WebMentions.AddWebMention(WebMention{
			Source: p.Permalink,
			Target: target,
			Status: WMStatusPending,
			Date:   p.Metadata.Date,
		})
	}
}
