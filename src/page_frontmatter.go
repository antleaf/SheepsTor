package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

type FrontMatter struct {
	Title       string    `yaml:"title"`
	Author      string    `yaml:"author"`
	Date        time.Time `yaml:"date"`
	Draft       bool      `yaml:"draft"`
	Description string    `yaml:"description"`
	Aliases     []string  `yaml:"aliases,omitempty"`
	Tags        []string  `yaml:"tags"`
	Collections []string  `yaml:"collections"`
	Editorial   []string  `yaml:"editorial"`
}

func (fm *FrontMatter) ToYaml() (string, error) {
	frontMatterString := ""
	frontMatterBytes, err := yaml.Marshal(fm)
	if err == nil {
		frontMatterString = "---\n" + string(frontMatterBytes) + "---\n\n"
	}
	return frontMatterString, err
}

func ParseTextForFrontMatterAndContent(rawPage string) (*FrontMatter, string, error) {
	var frontMatter FrontMatter
	var err error
	metadataAndContent := strings.SplitN(rawPage, "\n---\n", 2)
	if len(metadataAndContent) != 2 {
		logger.Debug("no metadata present")
		return nil, rawPage, errors.New("no metadata present")
	} else {
		err = yaml.Unmarshal([]byte(metadataAndContent[0]), &frontMatter)
		if err != nil {
			return nil, "", err
		}
	}
	return &frontMatter, metadataAndContent[1], err
}
