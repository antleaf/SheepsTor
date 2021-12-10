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
	Permalink string
	FilePath  string
	Slug      string
	Metadata  Frontmatter
	Content   string
}

type Frontmatter struct {
	Title       string    `yaml:"title"`
	Author      string    `yaml:"author"`
	Date        time.Time `yaml:"date"`
	Type        string    `yaml:"type"`
	Draft       bool      `yaml:"draft"`
	Description string    `yaml:"description"`
	Aliases     []string  `yaml:"aliases"`
	Tags        []string  `yaml:"tags"`
	Collections []string  `yaml:"collections"`
	Year        string    `yaml:"year"`
	Month       string    `yaml:"month"`
}

func (p *Page) GetFrontMatterAsYaml() (string, error) {
	frontMatterString := ""
	frontMatterBytes, err := yaml.Marshal(p.Metadata)
	if err == nil {
		frontMatterString = "---\n" + string(frontMatterBytes) + "---\n\n"
	}
	return frontMatterString, err
}

func (p *Page) WriteToFile(postFilePath string, withMetadata bool) error {
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
	err = ioutil.WriteFile(postFilePath, []byte(fullPostContent), os.ModePerm)
	return err
}

func (p *Page) ReadFromFile(postFilePath string) error {
	fileBytes, err := ioutil.ReadFile(postFilePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	slug := postFilePath[:len(postFilePath)-8]
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
	return err
}
