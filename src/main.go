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
	permalink := "https://www.paulwalk.net/2003/broadband-britain/"
	pp := w.PathProcessorSet.SelectPathProcessorForPermalink(permalink)
	logger.Infof("PP = %s", pp.Name)
	logger.Infof("Path = %s", w.GetPagePathForPermalink(permalink))
	//logger.Debugf("path = %s", w.GetPagePathForPermalink(permalink))
	//filePath := "posts/2003/broadband-britain/index.md"
	//filePath = "posts/2015/the-active-repository-pattern/index.md"
	//p, err := w.LoadPage(filePath)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Infof("Title of page = %s", p.Title)
	////logger.Infof("Webmention Count = %v", p.WebMentions.Count())
	//paths, _ := w.GetAllPageFilePaths()
	//for _, path := range paths {
	//	page, pageLoadErr := w.LoadPage(path)
	//	if pageLoadErr != nil {
	//		logger.Error(pageLoadErr.Error())
	//	}
	//	if page.WebMentions.Count() == 0 {
	//		logger.Infof("This page has %v webmentions: %s", page.WebMentions.Count(), page.FilePath)
	//	}
	//}
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
