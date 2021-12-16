package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var logger *zap.SugaredLogger
var SheepsTorConfig = Config{}
var systemReadyToAcceptUpdateRequests bool
var router chi.Router

func main() {
	systemReadyToAcceptUpdateRequests = false
	debugPtr := flag.Bool("debug", false, "-debug true|false")
	configFilePathPtr := flag.String("config", "./config/config.yaml", "-config <file_path>")
	updatePtr := flag.String("update", "", "-update all|<some_id>")
	flag.Parse()
	err := (&SheepsTorConfig).initialise(*debugPtr, *configFilePathPtr)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because SheepsTorConfig file not loaded from %s\n", *configFilePathPtr)
		os.Exit(1)
	}
	logger, err = ConfigureZapSugarLogger(SheepsTorConfig.DebugLogging)
	if SheepsTorConfig.DebugLogging {
		logger.Infof("Debugging enabled")
	}
	logger.Infof("WebRoot folder path set to: %s", SheepsTorConfig.WebRoot)
	logger.Infof("WebRoot folder path set to: %s", SheepsTorConfig.WebRoot)
	logger.Infof("Source Root folder path set to: %s", SheepsTorConfig.SourceRoot)
	SheepsTorConfig.configureWebsites()
	router = ConfigureRouter()
	//Scratch()
	if *updatePtr != "" {
		runAsCLIProcess(*updatePtr)
	} else {
		runAsHTTPProcess()
	}
}

func Scratch() {
	//w := SheepsTorConfig.getWebsiteByID("www.paulwalk.net")
	//for _, node := range w.SiteMap.Nodes {
	//	//page := node.LoadPage()
	//	//page.WriteToFile(true)
	//}
	os.Exit(1)
}

func runAsCLIProcess(sitesToUpdate string) {
	logger.Info(fmt.Sprintf("Running as CLI Process, updating website(s): '%s'...", sitesToUpdate))
	if sitesToUpdate == "all" {
		ProcessAllWebsites()
	} else {
		ProcessWebsite(SheepsTorConfig.getWebsiteByID(sitesToUpdate))
	}
}

func runAsHTTPProcess() {
	logger.Info(fmt.Sprintf("Running as HTTP Process on port %d", SheepsTorConfig.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", SheepsTorConfig.Port), router)
	if err != nil {
		logger.Error(err.Error())
	}
}
