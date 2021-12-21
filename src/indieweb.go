package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type IndieWeb struct {
	IndieAuthTokenEndpoint    string
	MicroPubMediaEndpoint     string
	IndieAuthId               string
	MediaUploadPath           string
	MediaUploadBaseURL        string
	DraftPosts                bool
	WebMentionIoWebhookSecret string
	MediaUploadURLRegex       *regexp.Regexp
}

func NewIndieWeb(iwConfig IndieWebConfig, baseURL, processorRoot string) IndieWeb {
	var iw = IndieWeb{
		IndieAuthTokenEndpoint: iwConfig.IndieAuthTokenEndpoint,
		MicroPubMediaEndpoint:  iwConfig.MicroPubMediaEndpoint,
		IndieAuthId:            iwConfig.IndieAuthId,
		DraftPosts:             iwConfig.DraftPosts,
		MediaUploadURLRegex:    nil,
	}
	iw.MediaUploadBaseURL = fmt.Sprintf("%s/media", baseURL)
	iw.MediaUploadPath = filepath.Join(processorRoot, "static", "media")
	iw.MediaUploadURLRegex = regexp.MustCompile("(" + iw.MediaUploadBaseURL + "/)([a-zA-Z0-9_-]+.[a-zA-Z]{2,4})")
	iw.WebMentionIoWebhookSecret = os.Getenv(iwConfig.WebMentionIoWebhookSecretEnvKey)
	return iw
}

type IndieAuthResult struct {
	Me       string `json:"me"`
	ClientId string `json:"client_id"`
	Scope    string `json:"scope"`
	Issue    int    `json:"issued_at"`
	Nonce    int    `json:"nonce"`
}

func CheckAccess(token, indieAuthTokenEndpoint, indieAuthMe string) (bool, error) {
	if token == "" {
		return false,
			errors.New("token string is empty")
	}
	// form the request to check the token
	client := &http.Client{}
	req, err := http.NewRequest("GET", indieAuthTokenEndpoint, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)
	// send the request
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	// parse the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	var indieAuthRes = IndieAuthResult{}
	err = json.Unmarshal(body, &indieAuthRes)
	if err != nil {
		return false, err
	}

	// verify results of the response
	if indieAuthRes.Me != indieAuthMe {
		return false, err
	}
	scopes := strings.Fields(indieAuthRes.Scope)
	postPresent := false
	for _, scope := range scopes {
		if scope == "post" || scope == "create" || scope == "update" {
			postPresent = true
			break
		}
	}
	if !postPresent {
		return false, errors.New("post is not present in the scope")
	}
	return true, nil
}
