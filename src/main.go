package main

import (
	"SheepsTor/src/sheepstor"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

var logger *zap.SugaredLogger
var config = Configuration{}
var router chi.Router

func main() {
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
	sheepstor.GitHubWebHookSecret = os.Getenv(config.GitHubWebHookSecretEnvKey)
	sheepstor.InitialiseRegistry(config.SourceRoot, config.WebRoot)
	for _, w := range config.WebsiteConfigs {
		website := sheepstor.NewWebsite(
			w.ID, w.ContentProcessor,
			w.ProcessorRootSubFolderPath,
			w.ContentRootSubFolderPath,
			config.SourceRoot,
			config.WebRoot,
			w.GitRepoConfig.CloneId,
			w.GitRepoConfig.RepoName,
			w.GitRepoConfig.BranchName,
		)
		sheepstor.Registry.Add(website)
	}
	logger.Infof("WebRoot folder path set to: %s", config.WebRoot)
	logger.Infof("Source Root folder path set to: %s", config.SourceRoot)
	//Scratch()
	if *updatePtr != "" {
		runAsCLIProcess(*updatePtr)
	} else {
		runAsHTTPProcess()
	}
}

func Scratch() {
	w := sheepstor.Registry.GetWebsiteByID("www.paulwalk.net")
	if w != nil {
		logger.Debugf("website ID = %s", w.ID)
	} else {
		logger.Debug("not found")
	}

	os.Exit(1)
}

func runAsCLIProcess(sitesToUpdate string) {
	logger.Info(fmt.Sprintf("Running as CLI Process, updating website(s): '%s'...", sitesToUpdate))
	if sitesToUpdate == "all" {
		processAllWebsites()
	} else {
		sheepstor.Registry.GetWebsiteByID(sitesToUpdate).ProcessWebsite()
	}
}

func runAsHTTPProcess() {
	router = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(10))
	//TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	//r.Handle("/_resources/assets/*", http.FileServer(http.FS(embeddedAssets)))
	router.Get("/", DefaultHandler)
	//r.Post("/comment", CommentPostHandler)
	router.Post("/update", sheepstor.GitHubWebHookHandler)
	logger.Info(fmt.Sprintf("Running as HTTP Process on port %d", config.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
	if err != nil {
		logger.Error(err.Error())
	}
}

func processWebsiteInSynchronousWorker(website *sheepstor.Website, wg *sync.WaitGroup) {
	err := website.ProcessWebsite()
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Infof("Processed website: '%s'", website.ID)
	}
	wg.Done()
}

func processAllWebsites() {
	var wg sync.WaitGroup
	for _, website := range sheepstor.Registry.WebSites {
		wg.Add(1)
		go processWebsiteInSynchronousWorker(website, &wg)
	}
	wg.Wait()
}

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello world"))
	resp.WriteHeader(http.StatusOK)
}
