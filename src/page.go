package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Page struct {
	Permalink   string
	FilePath    string
	Slug        string
	Metadata    Frontmatter
	Content     string
	Webmentions *WebmentionSet
	BaseURL     *string
}

type Frontmatter struct {
	Title       string    `yaml:"title"`
	Author      string    `yaml:"author"`
	Date        time.Time `yaml:"date"`
	Type        string    `yaml:"type"`
	Draft       bool      `yaml:"draft"`
	Description string    `yaml:"description"`
	Aliases     []string  `yaml:"aliases,omitempty"`
	Tags        []string  `yaml:"tags"`
	Collections []string  `yaml:"collections"`
}

func NewPage() Page {
	page := Page{}
	webmentionSet := make(WebmentionSet, 0)
	page.Webmentions = &webmentionSet
	return page
}

func (f *Frontmatter) Month() string {
	return f.Date.Format("01")
}

func (f *Frontmatter) Year() string {
	return f.Date.Format("2006")
}

func (p *Page) GetFrontMatterAsYaml() (string, error) {
	frontMatterString := ""
	frontMatterBytes, err := yaml.Marshal(p.Metadata)
	if err == nil {
		frontMatterString = "---\n" + string(frontMatterBytes) + "---\n\n"
	}
	return frontMatterString, err
}

func (p *Page) WriteToFile(withMetadata bool) error {
	var fullPostContent string
	var err error
	if withMetadata {
		fullPostContent, err = p.GetFrontMatterAsYaml()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	fullPostContent = fullPostContent + p.Content
	err = ioutil.WriteFile(p.FilePath, []byte(fullPostContent), os.ModePerm)
	p.GeneratePendingWebmentionsFromExtractedLinks()
	webmentionsErr := p.Webmentions.SaveToFile(filepath.Join(filepath.Dir(p.FilePath), "webmentions.yaml"))
	if webmentionsErr != nil {
		logger.Warn(err.Error())
	}
	return err
}

func (p *Page) ReadFromFile() error {
	var err error
	fileBytes, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	slug := p.FilePath[:len(p.FilePath)-8]
	slug = filepath.Base(slug)
	p.Slug = slug
	fileContent := string(fileBytes)
	fileContent = strings.TrimPrefix(fileContent, "---\n")
	metadataAndContent := strings.SplitN(fileContent, "\n---\n", 2)
	if len(metadataAndContent) != 2 {
		logger.Debug("no metadata present")
		p.Content = fileContent
	} else {
		metadata := metadataAndContent[0]
		err = yaml.Unmarshal([]byte(metadata), &p.Metadata)
		p.Content = metadataAndContent[1]
	}

	//READ WEBMENTIONS
	webmentionsErr := p.Webmentions.LoadFromFile(filepath.Join(filepath.Dir(p.FilePath), "webmentions.yaml"))
	if webmentionsErr != nil {
		//logger.Warn(webmentionsErr.Error())
	}
	return err
}

func (p *Page) GeneratePendingWebmentionsFromExtractedLinks() {
	targets := ExtractLinkURLs(p.Content, *p.BaseURL)
	for _, target := range targets {
		p.Webmentions.AddWebmention(Webmention{
			Source: p.Permalink,
			Target: target,
			Status: WMStatusPending,
			Date:   p.Metadata.Date,
		})
	}

}
