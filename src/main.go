package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

var logger *zap.SugaredLogger
var config = Configuration{}
var registry WebsiteRegistry
var systemReadyToAcceptUpdateRequests bool
var router chi.Router

func main() {
	systemReadyToAcceptUpdateRequests = false
	debugPtr := flag.Bool("debug", false, "-debug true|false")
	configFilePathPtr := flag.String("config", "./config/config.yaml", "-config <file_path>")
	updatePtr := flag.String("update", "", "-update all|<some_id>")
	flag.Parse()
	err := (&config).initialise(*debugPtr, *configFilePathPtr)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because SheepsTorConfig file not loaded from %s\n", *configFilePathPtr)
		os.Exit(1)
	}
	logger, err = ConfigureZapSugarLogger(config.DebugLogging)
	if config.DebugLogging {
		logger.Infof("Debugging enabled")
	}
	logger.Infof("WebRoot folder path set to: %s", config.WebRoot)
	logger.Infof("Source Root folder path set to: %s", config.SourceRoot)
	registry = NewRegistry(config.WebsiteConfigs, config.SourceRoot, config.WebRoot)
	router = ConfigureRouter()
	//Scratch()
	if *updatePtr != "" {
		runAsCLIProcess(*updatePtr)
	} else {
		runAsHTTPProcess()
	}
}

func Scratch() {
	w := registry.getWebsiteByID("www.paulwalk.net")
	//path := "posts/2008/why-i-suppose-i-ought-to-become-a-daily-mail-reader/index.md"
	//page, _ := w.LoadPage(path)
	//logger.Debugf("title: %s", page.Title)
	//for _, wm := range page.WebMentions.WebMentions {
	//	logger.Debugf("%s, %s, %s", wm.Status, wm.Source, wm.Target)
	//}

	paths, _ := w.GetAllPageFilePaths()
	for _, path := range paths {
		if strings.HasPrefix(path, "posts") {
			page, _ := w.LoadPage(path)
			//for _, wm := range page.WebMentions.WebMentions {
			//	if wm.Status == WMStatusPending {
			//		logger.Debug(wm.Source)
			//	}
			//}
			w.SavePage(page, page.FilePath)
		}
	}

	//w.DumpSiteMap( os.Stdout)
	os.Exit(1)
}

func runAsCLIProcess(sitesToUpdate string) {
	logger.Info(fmt.Sprintf("Running as CLI Process, updating website(s): '%s'...", sitesToUpdate))
	if sitesToUpdate == "all" {
		ProcessAllWebsites()
	} else {
		ProcessWebsite(registry.getWebsiteByID(sitesToUpdate))
	}
}

func runAsHTTPProcess() {
	logger.Info(fmt.Sprintf("Running as HTTP Process on port %d", config.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
	if err != nil {
		logger.Error(err.Error())
	}
}
