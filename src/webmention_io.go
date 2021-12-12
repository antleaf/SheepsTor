package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type WebmentionIOPayloadSourcePostAuthor struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
	URL   string `json:"url"`
}

type WebmentionIOPayloadSourcePost struct {
	Type       string                              `json:"type"`
	Author     WebmentionIOPayloadSourcePostAuthor `json:"author"`
	URL        string                              `json:"url"`
	Published  time.Time                           `json:"published"`
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
