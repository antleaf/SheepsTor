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

type WebMention struct {
	Source string    `csv:"source"`
	Target string    `csv:"target""`
	Status string    `csv:"status""`
	Date   time.Time `csv:"date""`
}

func (wm *WebMention) Send() {
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
		logger.Debugf("WebMention endpoint found for %s - attempting to send webmention", wm.Target)
		//logger.Debugf("Endpoint = %s", endpoint)
	}
	resp, err := client.SendWebmention(endpoint, wm.Source, wm.Target)
	if err != nil {
		wm.Status = WMStatusFailed
		logger.Errorf("WebMention failed to be sent for %s with error %s", wm.Target, err.Error())
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		wm.Status = WMStatusSent
		logger.Infof("WebMention sent successfully for %s with HTTP response %v", wm.Target, resp.StatusCode)
	} else {
		wm.Status = WMStatusFailed
		logger.Debugf("WebMention failed to be sent for %s with error %s", wm.Target, err.Error())
		return
	}
}

type WebMentionSet struct {
	WebMentions []WebMention
	FilePath    string
}

func NewWebMentionSet(filePath string) WebMentionSet {
	var wms = WebMentionSet{
		FilePath: filePath,
	}

	return wms
}

func (wms *WebMentionSet) LoadFromFile() error {
	f, err := os.Open(wms.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer f.Close()
	var tempWMSlice []WebMention
	err = gocsv.UnmarshalFile(f, &tempWMSlice)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	wms.WebMentions = make([]WebMention, 0)
	for _, wm := range tempWMSlice {
		wms.AddWebMention(wm)
	}
	return err
}

func (wms *WebMentionSet) SaveToFile() error {
	f, err := os.Create(wms.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer f.Close()
	err = gocsv.Marshal(wms.WebMentions, f)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

func (wms *WebMentionSet) AddWebMention(wm WebMention) {
	if wms.GetWebMentionBySourceAndTarget(wm.Source, wm.Target) == nil {
		wms.WebMentions = append(wms.WebMentions, wm)
	}
}

func (wms *WebMentionSet) Sort() {
	sort.Slice(wms.WebMentions, func(i, j int) bool {
		return (wms.WebMentions)[i].Date.Before((wms.WebMentions)[j].Date)
	})
}

func (wms *WebMentionSet) GetWebMentionBySourceAndTarget(source, target string) *WebMention {
	for _, wm := range wms.WebMentions {
		if (wm.Source == source) && (wm.Target == target) {
			return &wm
		}
	}
	return nil
}

func (wms *WebMentionSet) GetWebMentionBySource(source string) *WebMention {
	for _, wm := range wms.WebMentions {
		if wm.Source == source {
			return &wm
		}
	}
	return nil
}

func (wms *WebMentionSet) SendPendingWebMentions() error {
	var err error
	for _, wm := range wms.WebMentions {
		if wm.Status == WMStatusPending {
			logger.Debugf("Processing pending webmention from %s to %s", wm.Source, wm.Target)
			wm.Send()
		}
	}
	return err
}
