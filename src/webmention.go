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
	if wms.GetWebmentionBySourceAndTarget(wm.Source, wm.Target) == nil {
		*wms = append(*wms, &wm)
		logger.Debug("Added webmention")
	}
}

func (wms *WebmentionSet) Sort() {
	sort.Slice(wms, func(i, j int) bool {
		return (*wms)[i].Date.Before((*wms)[j].Date)
	})
}

func (wms *WebmentionSet) GetWebmentionBySourceAndTarget(source, target string) *Webmention {
	for _, wm := range *wms {
		if (wm.Source == source) && (wm.Target == target) {
			return wm
		}
	}
	return nil
}

func (wms *WebmentionSet) GetWebmentionBySource(source string) *Webmention {
	for _, wm := range *wms {
		if wm.Source == source {
			return wm
		}
	}
	return nil
}
