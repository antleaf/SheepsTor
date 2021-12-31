package main

import (
	"flag"
	"fmt"
	sheepstor "github.com/antleaf/sheepstor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

var logger *zap.SugaredLogger
var config = sheepstor.Configuration{}
var router chi.Router
var registry sheepstor.WebsiteRegistry

func main() {
	debugPtr := flag.Bool("debug", false, "-debug true|false")
	configFilePathPtr := flag.String("config", "./config.yaml", "-config <file_path>")
	updatePtr := flag.String("update", "", "-update all|<some_id>")
	flag.Parse()
	err := (&config).Initialise(*debugPtr, *configFilePathPtr)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because SheepsTorConfig file not loaded from %s\n", *configFilePathPtr)
		os.Exit(1)
	}
	logger, err = ConfigureZapSugarLogger(config.DebugLogging)
	if config.DebugLogging {
		logger.Infof("Debugging enabled")
	}
	sheepstor.SetLogger(logger)
	sheepstor.GitHubWebHookSecret = os.Getenv(config.GitHubWebHookSecretEnvKey)
	registry = sheepstor.NewRegistry(config.SourceRoot, config.WebRoot)
	for _, w := range config.WebsiteConfigs {
		website := sheepstor.NewWebsite(
			w.ID, w.ContentProcessor,
			w.ProcessorRootSubFolderPath,
			config.SourceRoot,
			config.WebRoot,
			w.GitRepoConfig.CloneId,
			w.GitRepoConfig.RepoName,
			w.GitRepoConfig.BranchName,
		)
		var websiteInterface sheepstor.WebsiteInterface
		websiteInterface = &website
		registry.Add(&websiteInterface)
	}
	logger.Infof("WebRoot folder path set to: %s", config.WebRoot)
	logger.Infof("Source Root folder path set to: %s", config.SourceRoot)
	if *updatePtr != "" {
		runAsCLIProcess(*updatePtr)
	} else {
		runAsHTTPProcess()
	}
}

func runAsCLIProcess(sitesToUpdate string) {
	logger.Info(fmt.Sprintf("Running as CLI Process, updating website(s): '%s'...", sitesToUpdate))
	if sitesToUpdate == "all" {
		processAllWebsites()
	} else {
		website := *registry.GetWebsiteByID(sitesToUpdate)
		err := website.ProvisionSources()
		if err != nil {
			logger.Error(err.Error())
			return
		}
		err = website.Build()
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func runAsHTTPProcess() {
	router = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(10))
	//TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	router.Get("/", DefaultHandler)
	router.Post("/update", registry.GitHubWebHookHandler)
	logger.Info(fmt.Sprintf("Running as HTTP Process on port %d", config.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
	if err != nil {
		logger.Error(err.Error())
	}
}

func processWebsiteInSynchronousWorker(websitePtr *sheepstor.WebsiteInterface, wg *sync.WaitGroup) {
	website := *websitePtr
	err := website.ProvisionSources()
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Infof("Provisioned sources for website: '%s'", website.GetID())
		err = website.Build()
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Infof("Built website: '%s'", website.GetID())
		}
	}
	wg.Done()
}

func processAllWebsites() {
	var wg sync.WaitGroup
	for _, website := range registry.WebSites {
		wg.Add(1)
		go processWebsiteInSynchronousWorker(website, &wg)
	}
	wg.Wait()
}

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello world"))
	resp.WriteHeader(http.StatusOK)
}
