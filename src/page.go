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
	Permalink     string
	FilePath      string
	Title         string
	Slug          string
	Published     time.Time
	Metadata      FrontMatter
	Content       string
	WebMentions   WebMentionSet
	PathProcessor *PathProcessor
}

func NewPage(filePath string, published time.Time, baseURL *string, pp *PathProcessor) *Page {
	page := Page{
		FilePath:      filePath,
		Published:     published,
		PathProcessor: pp,
	}
	page.Slug = filepath.Base(filepath.Dir(page.FilePath))
	page.Title = strings.ReplaceAll(page.Slug, "_", " ")
	page.Title = strings.Title(page.Title)
	page.ExtrapolatePermalink(baseURL)
	page.WebMentions = NewWebMentionSet()
	return &page
}

func (p *Page) WriteToFile(filePath string) error {
	var fullPostContent string
	var err error
	fullPostContent, err = p.Metadata.ToYaml()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	fullPostContent = fullPostContent + p.Content
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	err = ioutil.WriteFile(filePath, []byte(fullPostContent), os.ModePerm)
	p.WebMentions.SaveToFile(filepath.Join(filepath.Dir(filePath), "webmentions.csv"))
	return err
}

func (p *Page) ReadFromFile(filePath string) error {
	var err error
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = p.ReadFromString(string(fileBytes))
	if err != nil {
		return err
	}
	p.WebMentions.LoadFromFile(filepath.Join(filepath.Dir(filePath), "webmentions.csv"))
	return err
}

func (p *Page) ReadFromString(rawPage string) error {
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
	p.Content = content
	return err
}

func (p *Page) MoveMediaFromTempUploadToLocalFolderAndRewriteLinks(mediaUploadURLRegex *regexp.Regexp, mediaUploadPath, fullLocalFolderPath string) error {
	var err error
	for _, match := range mediaUploadURLRegex.FindAllStringSubmatch(p.Content, -1) {
		mediaFileName := match[2]
		logger.Debugf("Found media with file name %s", mediaFileName)
		oldFilePath := filepath.Join(mediaUploadPath, mediaFileName)
		newFilePath := filepath.Join(fullLocalFolderPath, mediaFileName)
		err = os.Rename(oldFilePath, newFilePath)
		if err != nil {
			logger.Warn(err.Error())
		}
	}
	p.Content = mediaUploadURLRegex.ReplaceAllString(p.Content, "./$2")
	return err
}

func (p *Page) ExtrapolatePermalink(baseURL *string) {
	permalink := p.PathProcessor.FolderRegex.ReplaceAllString(p.FilePath, p.PathProcessor.UrlGenerationPattern)
	permalink = strings.Replace(permalink, "{year}", p.Published.Format("2006"), -1)
	permalink = strings.Replace(permalink, "{month}", p.Published.Format("01"), -1)
	permalink = strings.Replace(permalink, "{slug}", p.Slug, -1)
	p.Permalink = *baseURL + "/" + permalink
}

func (p *Page) ExtractLinksForWebMentions(baseURL *string) {
	targets := ExtractLinkURLs(p.Content, *baseURL)
	for _, target := range targets {
		p.WebMentions.AddWebMention(WebMention{
			Source: p.Permalink,
			Target: target,
			Status: WMStatusPending,
			Date:   p.Metadata.Date,
		})
	}
}
