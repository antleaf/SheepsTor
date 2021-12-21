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

type WebMentionIOPayloadSourcePostAuthor struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
	URL   string `json:"url"`
}

type WebMentionIOPayloadSourcePost struct {
	Type       string                              `json:"type"`
	Author     WebMentionIOPayloadSourcePostAuthor `json:"author"`
	URL        string                              `json:"url"`
	Published  CustomTime                          `json:"published"`
	Name       string                              `json:"name"`
	RepostOf   string                              `json:"repost-of"`
	WmProperty string                              `json:"wm-property"`
}

type WebMentionIOPayload struct {
	Source  string                        `json:"source"`
	Target  string                        `json:"target"`
	Secret  string                        `json:"secret"`
	Deleted bool                          `json:"deleted,omitempty"`
	Post    WebMentionIOPayloadSourcePost `json:"post"`
}

func (w *WebMentionIOPayload) LoadAndValidate(payloadJson []byte, envNameForSecret string) (WebMention, error) {
	webMention := WebMention{}
	err := json.Unmarshal(payloadJson, w)
	if err != nil {
		logger.Error(err.Error())
		return webMention, err
	}
	if w.Secret != os.Getenv(envNameForSecret) {
		err = errors.New("secrets do not match - not authorised")
		logger.Error(err.Error())
		return webMention, err
	}
	webMention.Status = WMStatusReceived
	webMention.Source = w.Source
	webMention.Target = w.Target
	webMention.Date = time.Now()
	return webMention, err
}
