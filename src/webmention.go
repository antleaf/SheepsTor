package main

import (
	"github.com/gocarina/gocsv"
	"net/http"
	"os"
	"sort"
	"time"
	"willnorris.com/go/webmention"
)

const (
	WMStatusPending  string = "pending"
	WMStatusReceived        = "received"
	WMStatusSent            = "sent"
	WMStatusFailed          = "failed"
)

type Webmention struct {
	Source string    `csv:"source"`
	Target string    `csv:"target""`
	Status string    `csv:"status""`
	Date   time.Time `csv:"date""`
}

//type Webmention struct {
//	Source string `yaml:"source"`
//	Target string `yaml:"target"`
//	Status string `yaml:"status"`
//	Date   time.Time
//}

func (wm *Webmention) Send() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := webmention.New(httpClient)
	endpoint, err := client.DiscoverEndpoint(wm.Target)
	if err != nil {
		wm.Status = WMStatusFailed
		logger.Debugf("No endpoint found for %s", wm.Target)
		return
	} else {
		logger.Debugf("Webmention endpoint found for %s - attempting to send webmention", wm.Target)
		//logger.Debugf("Endpoint = %s", endpoint)
	}
	resp, err := client.SendWebmention(endpoint, wm.Source, wm.Target)
	if err != nil {
		wm.Status = WMStatusFailed
		logger.Errorf("Webmention failed to be sent for %s with error %s", wm.Target, err.Error())
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		wm.Status = WMStatusSent
		logger.Infof("Webmention sent successfully for %s with HTTP response %v", wm.Target, resp.StatusCode)
	} else {
		wm.Status = WMStatusFailed
		logger.Debugf("Webmention failed to be sent for %s with error %s", wm.Target, err.Error())
		return
	}
}

type WebmentionSet []*Webmention

func (wms *WebmentionSet) LoadFromFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer f.Close()
	tempWMSlice := []*Webmention{}
	err = gocsv.UnmarshalFile(f, &tempWMSlice)
	for _, wm := range tempWMSlice {
		wms.AddWebmention(*wm)
	}
	return err
}

func (wms *WebmentionSet) SaveToFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer f.Close()
	err = gocsv.Marshal(wms, f)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

func (wms *WebmentionSet) AddWebmention(wm Webmention) {
	if wms.GetWebmentionBySourceAndTarget(wm.Source, wm.Target) == nil {
		*wms = append(*wms, &wm)
		//logger.Debug("Added webmention")
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

func (wms *WebmentionSet) ProcessPending() error {
	var err error
	for _, wm := range *wms {
		if wm.Status == WMStatusPending {
			logger.Debugf("Processing pending webmention from %s to %s", wm.Source, wm.Target)
			wm.Send()
		}
	}
	return err
}
