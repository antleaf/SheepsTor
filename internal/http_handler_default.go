package internal

import "net/http"

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	err := Renderer.HTML(resp, http.StatusOK, "home", "Brave New World")
	if err != nil {
		Log.Error(err.Error())
	}
}
