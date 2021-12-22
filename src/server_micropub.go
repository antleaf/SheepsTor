package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func MicroPubAuthorisationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Authorising MicroPub request...")
		websiteID := chi.URLParam(r, "websiteID")
		website := registry.getWebsiteByID(websiteID)
		if website == nil {
			logger.Error(errors.New("website with ID " + websiteID + " not found"))
			http.Error(w, "not authorised", http.StatusUnauthorized)
			return
		}
		authorised, err := CheckAccess(r.Header.Get("Authorization"), website.IndieWeb.IndieAuthTokenEndpoint, website.IndieWeb.IndieAuthId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if authorised == false {
			logger.Error(errors.New("not authorised"))
			http.Error(w, "not authorised", http.StatusUnauthorized)
			return
		}
		logger.Info("Authorised!")
		ctx := context.WithValue(r.Context(), "websiteID", websiteID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MicroPubGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	website := registry.getWebsiteByID(r.Context().Value("websiteID").(string))
	if r.URL.Query().Get("q") == "config" {
		configResponse := MicroPubConfig{MediaEndpoint: website.IndieWeb.MicroPubMediaEndpoint}
		payload, _ := json.Marshal(configResponse)
		w.Write(payload)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func MicroPubMediaHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handling media upload...")
	website := registry.getWebsiteByID(r.Context().Value("websiteID").(string))
	err := r.ParseMultipartForm(20 << 20)
	//TODO: check what this 20 << 20 means above
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	mediaFile, header, err := r.FormFile("file")
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	mediaFilePath, mediaFileURL, err := website.StoreTempMediaFileAndReturnURL(mediaFile, header.Filename)
	logger.Debugf("Wrote file to %s", mediaFilePath)
	w.Header().Set("Location", mediaFileURL)
	w.WriteHeader(http.StatusCreated)
}

func MicroPubPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	website := registry.getWebsiteByID(r.Context().Value("websiteID").(string))
	payloadJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	entry := MicroPubPost{}
	err = json.Unmarshal(payloadJson, &entry)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	} else {
		entry.InitialisePublishedDate()
		filePath := fmt.Sprintf("posts/%s/%s/index.md", entry.Date.Format("2006"), entry.Properties.Name[0])
		pageIsNew := false
		page, pageLoadErr := website.LoadPage(filePath) //will return a new, blank page if not loaded from file
		if pageLoadErr != nil {
			pageIsNew = true
		}
		page.ReadFromString(entry.Properties.Content[0])
		err = website.SavePage(page, filePath)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		} else {
			logger.Debugf("Wrote page to file at %s", page.FilePath)
			page.MoveMediaFromTempUploadToLocalFolderAndRewriteLinks(website.IndieWeb.MediaUploadURLRegex, website.IndieWeb.MediaUploadPath, filepath.Dir(filepath.Join(website.ContentRoot, page.FilePath)))
			err = website.SavePage(page, filePath)
			w.Header().Set("Location", page.Permalink)
			if config.DisableGitCommitForDevelopment == false {
				err = website.CommitAndPush("Added or updated page on " + website.ID)
				if err != nil {
					logger.Error("Git Commit & Push failed " + err.Error())
					http.Error(w, "unable to commit changes to git", http.StatusBadRequest)
					return
				}
				err = website.Build()
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, "unable to rebuild website", http.StatusBadRequest)
					return
				}
				if pageIsNew {
					w.WriteHeader(http.StatusCreated)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
	//w.WriteHeader(http.StatusOK)
}
