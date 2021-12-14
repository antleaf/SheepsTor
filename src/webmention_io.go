package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const customTimeLayout = "2006-01-02T15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(customTimeLayout, s)
	return
}

type WebmentionIOPayloadSourcePostAuthor struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
	URL   string `json:"url"`
}

type WebmentionIOPayloadSourcePost struct {
	Type       string                              `json:"type"`
	Author     WebmentionIOPayloadSourcePostAuthor `json:"author"`
	URL        string                              `json:"url"`
	Published  CustomTime                          `json:"published"`
	Name       string                              `json:"name"`
	RepostOf   string                              `json:"repost-of"`
	WmProperty string                              `json:"wm-property"`
}

type WebmentionIOPayload struct {
	Source  string                        `json:"source"`
	Target  string                        `json:"target"`
	Secret  string                        `json:"secret"`
	Deleted bool                          `json:"deleted,omitempty"`
	Post    WebmentionIOPayloadSourcePost `json:"post"`
}

func (w *WebmentionIOPayload) LoadAndValidate(payloadJson []byte, envNameForSecret string) (Webmention, error) {
	webmention := Webmention{}
	//logger.Debug(string(payloadJson))
	err := json.Unmarshal(payloadJson, w)
	if err != nil {
		logger.Error(err.Error())
		return webmention, err
	}
	if w.Secret != os.Getenv(envNameForSecret) {
		err = errors.New("secrets do not match - not authorised")
		logger.Error(err.Error())
		return webmention, err
	}
	webmention.Status = WMStatusReceived
	webmention.Source = w.Source
	webmention.Target = w.Target
	webmention.Date = time.Now()
	return webmention, err
}
