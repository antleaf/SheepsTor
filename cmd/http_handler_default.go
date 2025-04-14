package cmd

import (
	"net/http"
)

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	err := Renderer.HTML(resp, http.StatusOK, "home", "This is SheepsTor")
	if err != nil {
		log.Error(err.Error())
	}
}
