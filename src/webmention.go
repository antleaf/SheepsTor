package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sort"
	"time"
)

const (
	WMStatusPending  string = "pending"
	WMStatusReceived        = "received"
	WMStatusSent            = "sent"
)

type Webmention struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Status string `yaml:"status"`
	Date   time.Time
}

type WebmentionSet []*Webmention

func (wms *WebmentionSet) LoadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = yaml.Unmarshal([]byte(data), &wms)
	return err
}

func (wms *WebmentionSet) SaveToFile(filePath string) error {
	data, err := yaml.Marshal(wms)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = ioutil.WriteFile(filePath, data, 0644)
	return err
}

func (wms *WebmentionSet) AddWebmention(wm Webmention) {
	*wms = append(*wms, &wm)
}

func (wms *WebmentionSet) Sort() {
	sort.Slice(wms, func(i, j int) bool {
		return (*wms)[i].Date.Before((*wms)[j].Date)
	})
}
