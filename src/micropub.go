package main

import (
	"time"
)

type MicroPubConfig struct {
	MediaEndpoint string `json:"media-endpoint"`
}

type MicroPubPostProperties struct {
	Name          []string  `json:"name"`
	Published     time.Time `json:"published"`
	Summary       []string  `json:"summary"`
	Category      []string  `json:"category"`
	Location      []string  `json:"location"`
	InReplyTo     []string  `json:"in-reply-to"`
	LikeOf        []string  `json:"like-of"`
	RepostOf      []string  `json:"repost-of"`
	Syndication   []string  `json:"syndication"`
	MPSyndicateTo []string  `json:"mp-syndicate-to"`
	Content       []string  `json:"content"`
}

type MicroPubPost struct {
	Type       []string               `json:"type"`
	Properties MicroPubPostProperties `json:"properties"`
	Date       time.Time              `json:"-"`
}

func (m *MicroPubPost) InitialisePublishedDate() {
	if m.Properties.Published.IsZero() {
		metadata, _, err := ParseTextForFrontMatterAndContent(m.Properties.Content[0])
		if err == nil {
			if metadata.Date.IsZero() == false {
				m.Date = metadata.Date
			} else {
				m.Date = time.Now()
			}
		}
	} else {
		m.Date = m.Properties.Published
	}
}
